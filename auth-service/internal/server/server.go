package server

import (
	"github.com/Elren44/go-auth/config"
	myhttp "github.com/Elren44/go-auth/internal/router"
	"net/http"
)

func New(app *config.App) *http.Server {
	mux := myhttp.NewRouter(app) // маршрутизатор

	return &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}
}
