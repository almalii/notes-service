package middlewares

import (
	"context"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"
	"notes-rew/internal/config"
	"notes-rew/internal/session"
	"time"
)

type SessionKey string

const sessionKey SessionKey = "session"

func SessionMiddleware(next http.Handler) http.Handler {
	cfg := config.InitConfig()
	sessionStore := sessions.NewCookieStore([]byte(cfg.Session))
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.SameSite = http.SameSiteStrictMode
	gob.Register(uuid.UUID{})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionStore.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = int(24 * time.Hour)

		ctx := context.WithValue(r.Context(), sessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SessMiddleware(next http.Handler) http.Handler {

	// Создание экземпляра вашего хранилища сессий
	sessionStore := session.NewRedisSessionStore("0.0.0.0:32768", "", 1)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлечение идентификатора сессии из заголовка
		sessionID := r.Header.Get("X-Session-ID")
		//if sessionID == "" {
		//	http.Error(w, "Session ID not provided", http.StatusBadRequest)
		//	return
		//}

		// Получение данных сессии из хранилища
		sessions, err := sessionStore.Get(r.Context(), sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Сохранение сессии в контексте запроса
		ctx := context.WithValue(r.Context(), sessionKey, sessions)

		// Вызов следующего обработчика в цепочке с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
