package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"notes-rew/internal/config"
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

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
