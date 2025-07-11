package router

import (
	"context"
	"github.com/Elren44/go-auth/config"
	"github.com/gookit/slog"
	"net/http"
	"strings"
	"time"
)

// Middleware-обёртка для логирования
func loggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Обёртка, чтобы перехватить статус-код
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		duration := time.Since(start)
		logger.Infof("%s %s %d %s", r.Method, r.URL.Path, rr.statusCode, duration)
	})
}

// Обёртка для записи статус-кода
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

type ctxKey string

const userIDKey ctxKey = "userID"

func JWTMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ParseAccessToken(tokenStr, app.Config.Secret)
			if err != nil {
				app.Logger.Warn("Invalid access token: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDKey).(int)
	return userID, ok
}
