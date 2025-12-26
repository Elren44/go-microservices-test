package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Elren44/go-gateway/config"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// extractToken extracts the JWT token from the request header or cookie
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return "", fmt.Errorf("invalid authorization header format")
		}
		return tokenString, nil
	}
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", fmt.Errorf("missing authorization header or access_token cookie")
	}
	return cookie.Value, nil
}

// attemptRefresh tries to refresh the access token using the refresh token
func attemptRefresh(r *http.Request, w http.ResponseWriter, cfg *config.GatewayConfig) (int, bool) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return 0, false
	}
	_, err = ParseRefreshToken(refreshCookie.Value, cfg.Secret)
	if err != nil {
		return 0, false
	}
	req, err := http.NewRequest("POST", cfg.AuthServiceURL+"/refresh", nil)
	if err != nil {
		return 0, false
	}
	req.AddCookie(refreshCookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return 0, false
	}
	var newUserID int
	var accessRefreshed bool
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "access_token" {
			uid, err := GetUserIDFromToken(cookie.Value, cfg.Secret)
			if err == nil {
				http.SetCookie(w, cookie)
				newUserID = uid
				accessRefreshed = true
			}
		} else if cookie.Name == "refresh_token" {
			http.SetCookie(w, cookie)
		}
	}
	resp.Body.Close()
	if accessRefreshed {
		return newUserID, true
	}
	return 0, false
}

// clearCookies clears the access and refresh token cookies
func clearCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})
}

func JWTMiddleware(cfg *config.GatewayConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userID, err := GetUserIDFromToken(token, cfg.Secret)
			log.Default().Println("JWTMiddleware userID:", userID, "err:", err)
			if err != nil {
				newUserID, ok := attemptRefresh(r, w, cfg)
				if !ok {
					clearCookies(w)
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
				userID = newUserID
			}

			// Add user ID to request header for downstream services
			r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
