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

func UserIdentity(tm *token_manager.TokenManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get(AuthorizationHeader)
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
				_, err = w.Write([]byte("invalid token"))
				if err != nil {
					return
				}
				return
			}

			parseUUID, err := uuid.Parse(userID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, err = w.Write([]byte("invalid user id"))
				if err != nil {
					return
				}
				return
			}

			ctx := context.WithValue(r.Context(), UserCtx, parseUUID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
