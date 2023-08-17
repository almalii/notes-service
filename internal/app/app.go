package app

import (
	"context"
	pb_auth_service "github.com/almalii/grpc-contracts/gen/go/auth_service/service/v1"
	pb_notes_service "github.com/almalii/grpc-contracts/gen/go/notes_service/service/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	authControllerGRPC "notes-rew/internal/auth_service/controller/grpc/v1"
	"notes-rew/internal/db/redis"
	"notes-rew/internal/middlewares"
	notesControllerGRPC "notes-rew/internal/notes_service/controller/grpc/v1"
	usersControllerGRPC "notes-rew/internal/users_service/controller/grpc/v1"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "notes-rew/docs"
	authController "notes-rew/internal/auth_service/controller/rest/handler"
	authService "notes-rew/internal/auth_service/service"
	authStorage "notes-rew/internal/auth_service/storage/postgres"
	authUsecase "notes-rew/internal/auth_service/usecase"
	"notes-rew/internal/config"
	"notes-rew/internal/db/migrations"
	"notes-rew/internal/db/postgres"
	"notes-rew/internal/hash"
	notesController "notes-rew/internal/notes_service/controller/rest/handler"
	notesService "notes-rew/internal/notes_service/service"
	notesStorage "notes-rew/internal/notes_service/storage/postgres"
	notesUsecase "notes-rew/internal/notes_service/usecase"
	"notes-rew/internal/token_manager"
	usersController "notes-rew/internal/users_service/controller/rest/handler"
	usersService "notes-rew/internal/users_service/service"
	usersStorage "notes-rew/internal/users_service/storage/postgres"
	usersUsecase "notes-rew/internal/users_service/usecase"
	"notes-rew/internal/validators"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
)

const (
	requestTimeout = 10 * time.Second
	contextTimeout = 10 * time.Second
	swaggerURL     = "http://localhost:8081/swagger/doc.json"
)

type grpcService struct {
	auth  pb_auth_service.AuthServiceServer
	users pb_users_service.UsersServiceServer
	notes pb_notes_service.NotesServiceServer
}

type App struct {
	protoService grpcService
	router       chi.Router
	mux          *runtime.ServeMux
	cfg          config.Config
	tokenManager *token_manager.TokenManager
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(requestTimeout))

	router.Mount("/debug", middleware.Profiler())

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerURL)))

	mux := runtime.NewServeMux()

	connectDB, err := postgres.ConnectionPostgresDB(ctx, cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %+v", err)
	}

	if err = migrations.UpMigrations(cfg); err != nil {
		logrus.Errorf("Failed to migrate: %+v", err)
	}

	connectRedis, err := redis.ConnectionRedisStorage(ctx, cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to Redis: %+v", err)
	}

	validation := validator.New()
	validators.RegisterCustomValidation(validation)

	tokenManager := token_manager.NewTokenManager(cfg.JwtSigning)

	hasher := hash.NewPasswordHasher(cfg.SaltHash)

	noteStorage := notesStorage.NewNoteStorage(connectDB)
	noteService := notesService.NewNoteService(noteStorage, connectRedis)
	noteUsecase := notesUsecase.NewNoteUsecase(noteService)
	noteController := notesController.NewNoteController(noteUsecase, validation, tokenManager)
	noteController.Register(router)

	noteControllerGRPC := notesControllerGRPC.NewNotesServer(
		noteUsecase,
		validation,
		pb_notes_service.UnimplementedNotesServiceServer{},
	)

	userStorage := usersStorage.NewPSQLUserStorage(connectDB)
	userService := usersService.NewUserService(userStorage)
	userUsecase := usersUsecase.NewUserUsecase(userService, hasher)
	userController := usersController.NewUserController(userUsecase, tokenManager, validation)
	userController.Register(router)

	userControllerGRPC := usersControllerGRPC.NewUsersServer(
		userUsecase,
		validation,
		pb_users_service.UnimplementedUsersServiceServer{},
	)

	authsStorage := authStorage.NewUserStorage(connectDB)
	authsService := authService.NewAuthService(authsStorage)
	authsUsecase := authUsecase.NewAuthUsecase(authsService, hasher, tokenManager)
	authsController := authController.NewAuthController(authsUsecase, validation)
	authsController.Register(router)

	authsControllerGRPC := authControllerGRPC.NewAuthServer(
		authsUsecase,
		validation,
		pb_auth_service.UnimplementedAuthServiceServer{},
	)

	return &App{
		router:       router,
		mux:          mux,
		cfg:          cfg,
		tokenManager: tokenManager,
		protoService: grpcService{
			auth:  authsControllerGRPC,
			users: userControllerGRPC,
			notes: noteControllerGRPC,
		},
	}
}

func (a *App) Start(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:           a.cfg.HTTPServer.Address,
		Handler:        a.router,
		ReadTimeout:    a.cfg.HTTPServer.ReadTimeout,
		WriteTimeout:   a.cfg.HTTPServer.WriteTimeout,
		MaxHeaderBytes: a.cfg.HTTPServer.MaxHeaderBytes,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	logrus.Println("HTTP server started on address:", a.cfg.HTTPServer.Address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	shutdownCtx, shutdown := context.WithTimeout(ctx, contextTimeout)
	defer shutdown()

	return httpServer.Shutdown(shutdownCtx)
}

func (a *App) StartGRPC() error {
	listener, err := net.Listen("tcp", a.cfg.GRPCServer.Address)
	if err != nil {
		logrus.Fatalf("Failed to listen: %+v", err)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		middlewares.UnaryTokenInterceptor(a.tokenManager)))

	pb_auth_service.RegisterAuthServiceServer(grpcServer, a.protoService.auth)
	pb_users_service.RegisterUsersServiceServer(grpcServer, a.protoService.users)
	pb_notes_service.RegisterNotesServiceServer(grpcServer, a.protoService.notes)

	reflection.Register(grpcServer)

	logrus.Println("gRPC server started on address:", a.cfg.GRPCServer.Address)

	if err = grpcServer.Serve(listener); err != nil {
		logrus.Fatalf("Failed to start gRPC server: %+v", err)
	}

	return nil
}

func (a *App) StartGateway(ctx context.Context) error {
	err := pb_auth_service.RegisterAuthServiceHandlerServer(ctx, a.mux, a.protoService.auth)
	if err != nil {
		return err
	}

	err = pb_users_service.RegisterUsersServiceHandlerServer(ctx, a.mux, a.protoService.users)
	if err != nil {
		return err
	}

	err = pb_notes_service.RegisterNotesServiceHandlerServer(ctx, a.mux, a.protoService.notes)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    a.cfg.GatewayServer.Address,
		Handler: middlewares.HttpInterceptor(a.tokenManager, a.mux),
	}

	logrus.Println("gRPC-Gateway server started on address:", a.cfg.GatewayServer.Address)

	if err = httpServer.ListenAndServe(); err != nil {
		logrus.Fatalf("Failed to start gRPC-Gateway server: %+v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	shutdownCtx, shutdown := context.WithTimeout(ctx, contextTimeout)
	defer shutdown()

	return httpServer.Shutdown(shutdownCtx)
}
