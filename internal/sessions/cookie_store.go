package sessions

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, sessionID, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(duration),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}

func ClearCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)
}

func GetSessionByCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func CheckCookieValue(w http.ResponseWriter, r *http.Request, cookieName string) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}

	if cookie.Value != "" {
		http.Error(w, "Cookie is not empty", http.StatusBadRequest)
		return true
	}

	return false
}
