package main

import (
	"fmt"
	"github.com/Elren44/go-auth/config"
	"github.com/Elren44/go-auth/internal/server"
	"log"
)

func main() {
	fmt.Println("Сервис авторизации")
	conf := config.NewAuthConfig()
	app := config.NewApp(conf)
	app.Logger.Info(conf)

	serv := server.New(app)
	log.Fatal(serv.ListenAndServe())

}
