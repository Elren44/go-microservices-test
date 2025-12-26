package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Elren44/go-gateway/config"
)

func NewRouter(cfg *config.GatewayConfig) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Proxy to auth-service
	authURL, _ := url.Parse(cfg.AuthServiceURL)
	authProxy := httputil.NewSingleHostReverseProxy(authURL)

	// Public auth routes (no JWT needed)
	mux.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/login" // Strip /auth prefix
		r.Method = "POST"
		authProxy.ServeHTTP(w, r)
	})
	mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/register"
		r.Method = "POST"
		authProxy.ServeHTTP(w, r)
	})
	mux.HandleFunc("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/refresh"
		authProxy.ServeHTTP(w, r)
	})
	mux.HandleFunc("POST /auth/logout", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/logout"
		r.Method = "POST"
		authProxy.ServeHTTP(w, r)
	})

	// Protected routes
	mux.Handle("/auth/me", JWTMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/me"
		authProxy.ServeHTTP(w, r)
	})))

	return mux
}
