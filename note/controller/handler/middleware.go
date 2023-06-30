package handler

import (
	"bookmarks/internal/config"
	"context"
	"github.com/gorilla/sessions"
	"net/http"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionStore := sessions.NewCookieStore([]byte(config.SessionKey()))
		session, err := sessionStore.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
