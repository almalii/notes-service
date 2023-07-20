package middlewares

import (
	"context"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"
	"notes-rew/internal/config"
	"time"
)

type SessionKey string

const sessionKey SessionKey = "sessions"

func SessionMiddleware(next http.Handler) http.Handler {
	cfg := config.InitConfig()
	sessionStore := sessions.NewCookieStore([]byte(cfg.Session))
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.SameSite = http.SameSiteStrictMode
	gob.Register(uuid.UUID{})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionStore.Get(r, "sessions")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = int(24 * time.Hour)

		ctx := context.WithValue(r.Context(), sessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
