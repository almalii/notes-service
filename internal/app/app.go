package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"net/http"
	authController "notes-rew/internal/auth/controller/rest/handler"
	authService "notes-rew/internal/auth/service"
	authStorage "notes-rew/internal/auth/storage"
	authUsecase "notes-rew/internal/auth/usecase"
	"notes-rew/internal/config"
	"notes-rew/internal/db/migrations"
	"notes-rew/internal/db/postgres"
	"notes-rew/internal/hash"
	notesController "notes-rew/internal/note/controller/handler"
	notesService "notes-rew/internal/note/service"
	notesStorage "notes-rew/internal/note/storage"
	notesUsecase "notes-rew/internal/note/usecase"
	usersController "notes-rew/internal/user/controller/handler"
	usersService "notes-rew/internal/user/service"
	usersStorage "notes-rew/internal/user/storage"
	usersUsecase "notes-rew/internal/user/usecase"
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
}

func NewApp() *App {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	logrus.SetFormatter(&logrus.JSONFormatter{})
	cfg := config.InitConfig()

	connectDB, err := postgres.ConnectionPostgresDB(context.Background(), cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %+v", err)
	}

	if err = migrations.UpMigrations(cfg); err != nil {
		logrus.Errorf("Failed to migrate: %+v", err)
	}

	validation := validator.New()
	validators.RegisterCustomValidation(validation)
	hasher := hash.NewPasswordHasher(cfg.Salt)

	noteStorage := notesStorage.NewNoteStorage(connectDB)
	noteService := notesService.NewNoteService(noteStorage)
	noteUsecase := notesUsecase.NewNoteUsecase(noteService)
	noteController := notesController.NewNoteController(noteUsecase, validation)
	noteController.Register(router)

	userStorage := usersStorage.NewPSQLUserStorage(connectDB)
	userService := usersService.NewUserService(userStorage)
	userUsecase := usersUsecase.NewUserUsecase(userService, hasher)
	userController := usersController.NewUserController(userUsecase, validation)
	userController.Register(router)

	authsStorage := authStorage.NewUserStorage(connectDB)
	authsService := authService.NewAuthService(authsStorage)
	authsUsecase := authUsecase.NewAuthUsecase(authsService, hasher)
	authsController := authController.NewAuthController(authsUsecase, validation)
	authsController.Register(router)

	return &App{router: router}
}

func (a *App) Start(port string) error {
	httpServer := &http.Server{
		Addr:           ":" + port,
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

	logrus.Println("server started on port " + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return httpServer.Shutdown(ctx)
}
