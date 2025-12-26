package router

import (
	"net/http"
	"time"

	"github.com/Elren44/go-auth/config"
	"github.com/gookit/slog"
)

func loggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("Request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

func NewRouter(app *config.App) http.Handler {
	mux := http.NewServeMux()

	//публичные роуты
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, app)
	})

	mux.HandleFunc("POST /refresh", func(w http.ResponseWriter, r *http.Request) {
		refreshHandler(w, r, app)
	})

	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(w, r, app)
	})
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler(w, r, app)
	})

	//приватные роуты

	mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		meHandler(w, r, app)
	})

	return loggingMiddleware(mux, app.Logger)
}
