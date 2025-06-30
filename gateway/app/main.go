package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("API Шлюз!")

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
