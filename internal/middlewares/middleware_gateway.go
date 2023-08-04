package middlewares

import (
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"notes-rew/internal/token_manager"
	"strings"
)

const (
	requestEndpoint1 = "register"
	requestEndpoint2 = "login"
)

func isAuthRequest(req string) bool {
	return strings.Contains(req, requestEndpoint1) || strings.Contains(req, requestEndpoint2)
}

func HttpInterceptor(tm *token_manager.TokenManager, next *runtime.ServeMux) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if isAuthRequest(req.URL.Path) {
			next.ServeHTTP(w, req)
			return
		}

		authHeader := req.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("empty auth header"))
			if err != nil {
				return
			}
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("invalid auth header"))
			if err != nil {
				return
			}
			return
		}

		userID, err := tm.ParseToken(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parseUUID, err := uuid.Parse(userID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), UserCtx, parseUUID)

		next.ServeHTTP(w, req.WithContext(ctx))
	}
}
