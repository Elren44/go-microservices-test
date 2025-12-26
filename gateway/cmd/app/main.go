package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Elren44/go-gateway/config"
	"github.com/Elren44/go-gateway/internal/router"
)

func main() {
	fmt.Println("API Шлюз!")
	cfg := config.NewGatewayConfig()
	fmt.Println(cfg)

	mux := router.NewRouter(cfg)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
