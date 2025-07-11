package router

import (
	"encoding/json"
	"fmt"
	"github.com/Elren44/go-auth/config"
	"net/http"
	"time"
)

func loginHandler(w http.ResponseWriter, r *http.Request, app *config.App) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		app.Logger.Error("Bad request: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Примитивная авторизация
	if creds.Username != "admin" || creds.Password != "password" {
		app.Logger.Warnf("Unauthorized login attempt: %s", creds.Username)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := 123 // <- вот здесь ты его задаёшь (можно получить из базы в будущем)

	accessToken, _ := GenerateAccessToken(userID, app.Config.Secret)
	refreshToken, _ := GenerateRefreshToken(userID, app.Config.Secret)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
func refreshHandler(w http.ResponseWriter, r *http.Request, app *config.App) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	// ✅ Парсим refresh token
	claims, err := ParseRefreshToken(cookie.Value, app.Config.Secret)
	if err != nil {
		app.Logger.Warn("Invalid refresh token: ", err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// ✅ Генерируем новый access token
	accessToken, err := GenerateAccessToken(claims.UserID, app.Config.Secret)
	if err != nil {
		app.Logger.Error("Access token generation failed: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// (опционально) обновляем refresh токен (rotation strategy)
	// newRefresh, _ := auth.GenerateRefreshToken(claims.UserID, app.Config.Secret)
	// http.SetCookie(w, &http.Cookie{...})

	// ✅ Возвращаем новый access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func meHandler(w http.ResponseWriter, r *http.Request, app *config.App) {
	userID, ok := GetUserID(r)
	if !ok {
		app.Logger.Warn("Unauthorized request")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	app.Logger.Info("User ID:", userID)
	fmt.Fprintf(w, "Your user ID: %d", userID)
}

func registerHandler(w http.ResponseWriter, r *http.Request, app *config.App) {
	w.WriteHeader(http.StatusOK)
	app.Logger.Info("register")
	w.Write([]byte("register"))
}

func logoutHandler(w http.ResponseWriter, r *http.Request, app *config.App) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
	w.WriteHeader(http.StatusOK)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
