package handler

import (
	"context"
	"github.com/gorilla/sessions"
	"net/http"
	"notes-rew/internal/config"
	"time"
)

func SessionMiddleware(next http.Handler) http.Handler {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey()))
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.SameSite = http.SameSiteStrictMode

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionStore.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = int(24 * time.Hour)

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
