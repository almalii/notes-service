package app

import (
	"context"
	pb_auth_service "github.com/almalii/grpc-contracts/gen/go/auth_service/service/v1"
	pb_notes_service "github.com/almalii/grpc-contracts/gen/go/notes_service/service/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	authControllerGRPC "notes-rew/internal/auth_service/controller/grpc/v1"
	authController "notes-rew/internal/auth_service/controller/rest/handler"
	authService "notes-rew/internal/auth_service/service"
	authStorage "notes-rew/internal/auth_service/storage"
	authUsecase "notes-rew/internal/auth_service/usecase"
	"notes-rew/internal/config"
	"notes-rew/internal/db/migrations"
	"notes-rew/internal/db/postgres"
	"notes-rew/internal/hash"
	notesControllerGRPC "notes-rew/internal/notes_service/controller/grpc/v1"
	notesController "notes-rew/internal/notes_service/controller/rest/handler"
	notesService "notes-rew/internal/notes_service/service"
	notesStorage "notes-rew/internal/notes_service/storage"
	notesUsecase "notes-rew/internal/notes_service/usecase"
	"notes-rew/internal/sessions"
	usersControllerGRPC "notes-rew/internal/users_service/controller/grpc/v1"
	usersController "notes-rew/internal/users_service/controller/rest/handler"
	usersService "notes-rew/internal/users_service/service"
	usersStorage "notes-rew/internal/users_service/storage"
	usersUsecase "notes-rew/internal/users_service/usecase"
	"notes-rew/internal/validators"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type AppGRPC struct {
	srv *grpc.Server
	ctx context.Context
	cfg config.Config
}

type App struct {
	router chi.Router
	ctx    context.Context
	cfg    config.Config
}

func NewApp(ctx context.Context, cfg config.Config) *App {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	connectDB, err := postgres.ConnectionPostgresDB(ctx, cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %+v", err)
	}

	if err = migrations.UpMigrations(cfg); err != nil {
		logrus.Errorf("Failed to migrate: %+v", err)
	}

	validation := validator.New()
	validators.RegisterCustomValidation(validation)

	sessionStore := sessions.NewRedisSessionStore("localhost:32768", "", 0)

	hasher := hash.NewPasswordHasher(cfg.Salt)

	noteStorage := notesStorage.NewNoteStorage(connectDB)
	noteService := notesService.NewNoteService(noteStorage)
	noteUsecase := notesUsecase.NewNoteUsecase(noteService)
	noteController := notesController.NewNoteController(noteUsecase, validation, sessionStore)
	noteController.Register(router)

	userStorage := usersStorage.NewPSQLUserStorage(connectDB)
	userService := usersService.NewUserService(userStorage)
	userUsecase := usersUsecase.NewUserUsecase(userService, hasher)
	userController := usersController.NewUserController(userUsecase, validation, sessionStore)
	userController.Register(router)

	authsStorage := authStorage.NewUserStorage(connectDB)
	authsService := authService.NewAuthService(authsStorage)
	authsUsecase := authUsecase.NewAuthUsecase(authsService, hasher)
	authsController := authController.NewAuthController(authsUsecase, validation, sessionStore)
	authsController.Register(router)

	return &App{
		router: router,
		ctx:    ctx,
		cfg:    cfg,
	}
}

func NewAppGRPC(ctx context.Context, cfg config.Config) *AppGRPC {
	return &AppGRPC{
		srv: grpc.NewServer(),
		ctx: ctx,
		cfg: cfg,
	}
}

func (ap *AppGRPC) StartGRPC() error {
	listener, err := net.Listen("tcp", ap.cfg.GRPCServer.Address)
	if err != nil {
		logrus.Fatalf("Failed to listen: %+v", err)
	}

	pb_auth_service.RegisterAuthServiceServer(
		ap.srv,
		authControllerGRPC.NewAuthServer(pb_auth_service.UnimplementedAuthServiceServer{}),
	)
	pb_users_service.RegisterUsersServiceServer(
		ap.srv,
		usersControllerGRPC.NewUsersServer(pb_users_service.UnimplementedUsersServiceServer{}),
	)
	pb_notes_service.RegisterNotesServiceServer(
		ap.srv,
		notesControllerGRPC.NewNotesServer(pb_notes_service.UnimplementedNotesServiceServer{}),
	)

	reflection.Register(ap.srv)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err = pb_auth_service.RegisterAuthServiceHandlerFromEndpoint(
		ap.ctx,
		mux,
		ap.cfg.GRPCServer.Address,
		opts,
	); err != nil {
		logrus.Fatalf("Failed to register auth service v1: %+v", err)
	}

	if err = pb_users_service.RegisterUsersServiceHandlerFromEndpoint(
		ap.ctx,
		mux,
		ap.cfg.GRPCServer.Address,
		opts,
	); err != nil {
		logrus.Fatalf("Failed to register users service v1: %+v", err)
	}

	if err = pb_notes_service.RegisterNotesServiceHandlerFromEndpoint(
		ap.ctx,
		mux,
		ap.cfg.GRPCServer.Address,
		opts,
	); err != nil {
		logrus.Fatalf("Failed to register notes service v1: %+v", err)
	}

	g, _ := errgroup.WithContext(ap.ctx)

	g.Go(func() error {
		log.Println("grpc server started on address:", ap.cfg.GRPCServer.Address)
		return ap.srv.Serve(listener)
	})

	g.Go(func() error {
		log.Println("grpc-gateway server started on address:", ap.cfg.GRPCServer.GateWayAddress)
		return http.ListenAndServe(ap.cfg.GRPCServer.GateWayAddress, mux)
	})

	return g.Wait()
}

func (a *App) StartHTTP() error {
	httpServer := &http.Server{
		Addr:           a.cfg.HTTPServer.Address,
		Handler:        a.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	log.Println("server started on address:", a.cfg.HTTPServer.Address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	shutdownCtx, shutdown := context.WithTimeout(a.ctx, 5*time.Second)
	defer shutdown()

	return httpServer.Shutdown(shutdownCtx)
}
