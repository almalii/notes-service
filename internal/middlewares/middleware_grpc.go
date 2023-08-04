package middlewares

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"notes-rew/internal/token_manager"
	"strings"
)

const grpcService = "AuthService"

func isAuthMethod(info string) bool {
	return strings.Contains(info, grpcService)
}

func UnaryTokenInterceptor(tm *token_manager.TokenManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if isAuthMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get(AuthorizationHeader)
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

		ctx = context.WithValue(ctx, UserCtx, parseUUID)

		return handler(ctx, req)
	}
}
