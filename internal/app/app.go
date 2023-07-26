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

type Deps struct {
	AuthHandler  pb_auth_service.AuthServiceServer
	UsersHandler pb_users_service.UserServiceServer
	NotesHandler pb_notes_service.NotesServiceServer
}

type AppGRPC struct {
	Deps
	srv *grpc.Server
	ctx context.Context
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

func NewAppGRPC(ctx context.Context) *AppGRPC {
	return &AppGRPC{
		srv: grpc.NewServer(),
		ctx: ctx,
	}
}

func (ap *AppGRPC) StartGRPC() error {

	pb_auth_service.RegisterAuthServiceServer(ap.srv, authControllerGRPC.NewAuthServer(pb_auth_service.UnimplementedAuthServiceServer{}))
	pb_users_service.RegisterUserServiceServer(ap.srv, usersControllerGRPC.NewUsersServer(pb_users_service.UnimplementedUserServiceServer{}))
	pb_notes_service.RegisterNotesServiceServer(ap.srv, notesControllerGRPC.NewNotesServer(pb_notes_service.UnimplementedNotesServiceServer{}))

	reflection.Register(ap.srv)

	listener, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		logrus.Fatalf("Failed to listen: %+v", err)
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb_auth_service.RegisterAuthServiceHandlerFromEndpoint(ap.ctx, mux, ":8090", opts)
	if err != nil {
		logrus.Fatalf("Failed to register auth service v1: %+v", err)
	}
	err = pb_users_service.RegisterUserServiceHandlerFromEndpoint(ap.ctx, mux, ":8090", opts)
	if err != nil {
		logrus.Fatalf("Failed to register users service v1: %+v", err)
	}
	err = pb_notes_service.RegisterNotesServiceHandlerFromEndpoint(ap.ctx, mux, ":8090", opts)
	if err != nil {
		logrus.Fatalf("Failed to register notes service v1: %+v", err)
	}

	g, _ := errgroup.WithContext(ap.ctx)

	g.Go(func() error {
		logrus.Println("grpc server started on address:", "localhost:8090")
		return ap.srv.Serve(listener)
	})

	g.Go(func() error {
		logrus.Println("grpc-gateway server started on address:", "localhost:8080")
		return http.ListenAndServe("localhost:8080", mux)
	})

	return g.Wait()
}

func (a *App) Start() error {
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
