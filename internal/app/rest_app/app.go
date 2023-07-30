package rest_app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
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
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

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

	tokenManager := token_manager.NewTokenManager(cfg.JwtSigning)
	hasher := hash.NewPasswordHasher(cfg.Salt)

	noteStorage := notesStorage.NewNoteStorage(connectDB)
	noteService := notesService.NewNoteService(noteStorage)
	noteUsecase := notesUsecase.NewNoteUsecase(noteService, validation)
	noteController := notesController.NewNoteController(noteUsecase, tokenManager)
	noteController.Register(router)

	userStorage := usersStorage.NewPSQLUserStorage(connectDB)
	userService := usersService.NewUserService(userStorage)
	userUsecase := usersUsecase.NewUserUsecase(userService, hasher, validation)
	userController := usersController.NewUserController(userUsecase, tokenManager)
	userController.Register(router)

	authsStorage := authStorage.NewUserStorage(connectDB)
	authsService := authService.NewAuthService(authsStorage)
	authsUsecase := authUsecase.NewAuthUsecase(authsService, hasher, tokenManager, validation)
	authsController := authController.NewAuthController(authsUsecase)
	authsController.Register(router)

	return &App{
		router: router,
		ctx:    ctx,
		cfg:    cfg,
	}
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
