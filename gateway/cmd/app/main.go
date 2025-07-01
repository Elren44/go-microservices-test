package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Elren44/go-gateway/config"
)

func main() {
	fmt.Println("API Шлюз!")
	config := config.NewAuthConfig()
	fmt.Println(config)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
