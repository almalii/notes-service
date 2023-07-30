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

func StreamTokenInterceptor(tm token_manager.TokenManager) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Извлекаем токен из метаданных.
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		token, ok := md["authorization"]
		if !ok || len(token) == 0 {
			return status.Errorf(codes.Unauthenticated, "missing token")
		}

		// Проверяем токен.
		claims, err := tm.ParseToken(token[0])
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Вносим информацию о пользователя в контекст запроса.
		ctx := context.WithValue(ss.Context(), "userID", claims)

		// Выполняем обработку запроса.
		return handler(srv, &wrappedServerStream{ServerStream: ss, ctx: ctx})
	}
}

// wrappedServerStream оборачивает grpc.ServerStream для передачи контекста.
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func HttpInterceptor() runtime.ServeMuxOption {
	httpServerOptions := runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		return nil
	})
	return httpServerOptions
}
