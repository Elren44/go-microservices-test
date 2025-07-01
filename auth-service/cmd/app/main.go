package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Elren44/go-auth/config"
)

func main() {
	fmt.Println("Сервис авторизации")
	config := config.NewAuthConfig()
	logger := log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	logger.Println(config)
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // 204 No Content
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Запрос: метод=%s, путь=%s, от=%s\n", r.Method, r.URL.Path, r.RemoteAddr)
		_, err := fmt.Fprintf(w, "Hello World")
		if err != nil {
			return
		}
	})

	log.Fatal(http.ListenAndServe(":4000", mux))
}
