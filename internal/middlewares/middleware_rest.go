package middlewares

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"notes-rew/internal/token_manager"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	UserCtx             = "userID"
)

func UserIdentity(tm token_manager.TokenManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get(AuthorizationHeader)
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("empty auth header"))
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid auth header"))
				return
			}

			userID, err := tm.ParseToken(headerParts[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid token"))
				return
			}

			parseUUID, err := uuid.Parse(userID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid user id"))
				return
			}

			ctx := context.WithValue(r.Context(), UserCtx, parseUUID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
