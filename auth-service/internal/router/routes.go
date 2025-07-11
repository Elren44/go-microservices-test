package router

import (
	"github.com/Elren44/go-auth/config"
	"net/http"
)

func NewRouter(app *config.App) http.Handler {
	mux := http.NewServeMux()

	//публичные роуты
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, app)
	})

	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		refreshHandler(w, r, app)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(w, r, app)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler(w, r, app)
	})

	//приватные роуты

	mux.Handle("/me", JWTMiddleware(app)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meHandler(w, r, app)
	})))

	return loggingMiddleware(mux, app.Logger)
}
