package grpc_app

import (
	"context"
	pb_auth_service "github.com/almalii/grpc-contracts/gen/go/auth_service/service/v1"
	pb_notes_service "github.com/almalii/grpc-contracts/gen/go/notes_service/service/v1"
	pb_users_service "github.com/almalii/grpc-contracts/gen/go/users_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	authControllerGRPC "notes-rew/internal/auth_service/controller/grpc/v1"
	authService "notes-rew/internal/auth_service/service"
	authStorage "notes-rew/internal/auth_service/storage"
	authUsecase "notes-rew/internal/auth_service/usecase"
	"notes-rew/internal/config"
	"notes-rew/internal/db/migrations"
	"notes-rew/internal/db/postgres"
	"notes-rew/internal/hash"
	notesControllerGRPC "notes-rew/internal/notes_service/controller/grpc/v1"
	notesService "notes-rew/internal/notes_service/service"
	notesStorage "notes-rew/internal/notes_service/storage"
	notesUsecase "notes-rew/internal/notes_service/usecase"
	usersControllerGRPC "notes-rew/internal/users_service/controller/grpc/v1"
	usersService "notes-rew/internal/users_service/service"
	usersStorage "notes-rew/internal/users_service/storage"
	usersUsecase "notes-rew/internal/users_service/usecase"
	"notes-rew/internal/validators"
)

type grpcService struct {
	auth  pb_auth_service.AuthServiceServer
	users pb_users_service.UsersServiceServer
	notes pb_notes_service.NotesServiceServer
}

type AppGRPC struct {
	grpcServer    *grpc.Server
	httpServer    *http.Server
	serviceServer grpcService
	ctx           context.Context
	cfg           config.Config
}

func NewAppGRPC(ctx context.Context, cfg config.Config) *AppGRPC {
	connectDB, err := postgres.ConnectionPostgresDB(ctx, cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %+v", err)
	}

	if err = migrations.UpMigrations(cfg); err != nil {
		logrus.Errorf("Failed to migrate: %+v", err)
	}

	validation := validator.New()
	validators.RegisterCustomValidation(validation)

	hasher := hash.NewPasswordHasher(cfg.Salt)

	authsStorage := authStorage.NewUserStorage(connectDB)
	authsService := authService.NewAuthService(authsStorage)
	authsUsecase := authUsecase.NewAuthUsecase(authsService, hasher)
	authsControllerGRPC := authControllerGRPC.NewAuthServer(
		authsUsecase,
		validation,
		pb_auth_service.UnimplementedAuthServiceServer{},
	)

	userStorage := usersStorage.NewPSQLUserStorage(connectDB)
	userService := usersService.NewUserService(userStorage)
	userUsecase := usersUsecase.NewUserUsecase(userService, hasher)
	userControllerGRPC := usersControllerGRPC.NewUsersServer(
		userUsecase,
		validation,
		pb_users_service.UnimplementedUsersServiceServer{},
	)

	noteStorage := notesStorage.NewNoteStorage(connectDB)
	noteService := notesService.NewNoteService(noteStorage)
	noteUsecase := notesUsecase.NewNoteUsecase(noteService)
	noteControllerGRPC := notesControllerGRPC.NewNotesServer(
		noteUsecase,
		validation,
		pb_notes_service.UnimplementedNotesServiceServer{},
	)

	return &AppGRPC{
		grpcServer: grpc.NewServer(),
		ctx:        ctx,
		cfg:        cfg,
		serviceServer: grpcService{
			auth:  authsControllerGRPC,
			users: userControllerGRPC,
			notes: noteControllerGRPC,
		},
	}
}

func (ap *AppGRPC) StartGRPC() error {
	listener, err := net.Listen("tcp", ap.cfg.GRPCServer.Address)
	if err != nil {
		logrus.Fatalf("Failed to listen: %+v", err)
	}

	var serverOptions []grpc.ServerOption

	ap.grpcServer = grpc.NewServer(serverOptions...)

	pb_auth_service.RegisterAuthServiceServer(ap.grpcServer, ap.serviceServer.auth)
	pb_users_service.RegisterUsersServiceServer(ap.grpcServer, ap.serviceServer.users)
	pb_notes_service.RegisterNotesServiceServer(ap.grpcServer, ap.serviceServer.notes)

	reflection.Register(ap.grpcServer)

	return ap.grpcServer.Serve(listener)
}

func (ap *AppGRPC) StartHTTP() error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb_auth_service.RegisterAuthServiceHandlerFromEndpoint(
		ap.ctx,
		mux,
		ap.cfg.GRPCServer.Address,
		opts,
	); err != nil {
		logrus.Fatalf("Failed to register auth service v1: %+v", err)
	}

	if err := pb_users_service.RegisterUsersServiceHandlerFromEndpoint(
		ap.ctx,
		mux,
		ap.cfg.GRPCServer.Address,
		opts,
	); err != nil {
		logrus.Fatalf("Failed to register users service v1: %+v", err)
	}

	if err := pb_notes_service.RegisterNotesServiceHandlerFromEndpoint(
		ap.ctx,
		mux,
		ap.cfg.GRPCServer.Address,
		opts,
	); err != nil {
		logrus.Fatalf("Failed to register notes service v1: %+v", err)
	}

	return http.ListenAndServe(ap.cfg.HTTPServer.Address, mux)

}
