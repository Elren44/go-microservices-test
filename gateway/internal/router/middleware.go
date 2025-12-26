package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Elren44/go-gateway/config"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWTMiddleware(cfg *config.GatewayConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			// First, try to get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString == authHeader {
					http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
					return
				}
			} else {
				// If no header, try to get token from cookie
				cookie, err := r.Cookie("access_token")
				if err != nil {
					http.Error(w, "Missing Authorization header or access_token cookie", http.StatusUnauthorized)
					return
				}
				tokenString = cookie.Value
			}

			userID, err := GetUserIDFromToken(tokenString, cfg.Secret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user ID to request header for downstream services
			r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
