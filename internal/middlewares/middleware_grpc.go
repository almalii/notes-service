package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"notes-rew/internal/token_manager"
	"strings"
)

type validatable interface {
	Validate() error
}

func UnaryTokenInterceptor(tm token_manager.TokenManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Ваша логика проверки токена здесь.
		// Вы можете получить токен из метаданных вызова и проверить его с помощью вашего TokenManager.

		// Пример:
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) != 1 {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header")
		}

		headerParts := strings.Split(authHeader[0], " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header")
		}

		userID, err := tm.ParseToken(headerParts[1])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		parseUUID, err := uuid.Parse(userID)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user id")
		}

		// Создайте контекст с информацией о пользователе и передайте его дальше обработчику.
		ctx = context.WithValue(ctx, UserCtx, parseUUID)

		// Вызовите обработчик для продолжения обработки запроса.
		return handler(ctx, req)
	}
}

func GrpcInterceptor() grpc.ServerOption {
	grpcServerOptions := grpc.UnaryInterceptor(func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if v, ok := req.(validatable); ok {
			err := v.Validate()
			if err != nil {
				return nil, err
			}
		}
		resp, err = handler(ctx, req)
		return handler(ctx, req)
	})
	return grpcServerOptions
}

func HttpInterceptor() runtime.ServeMuxOption {
	httpServerOptions := runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		return nil
	})
	return httpServerOptions
}

//func (m *MicroserviceServer) getUserIdFromToken(ctx context.Context) (string, error) {
//	md, _ := metadata.FromIncomingContext(ctx)
//	token := md.Get("Authorization")
//	if token == nil {
//		return "", status.Errorf(codes.PermissionDenied, "user isn't authorized")
//	}
//
//	userID, err := m.tokenManager.Parse(token[0])
//	if err != nil {
//		return "", err
//	}
//	return *userID, nil
//}
