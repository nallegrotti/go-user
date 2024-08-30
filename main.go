package main

import (
	"fmt"
	"log"
	"net/http"

	"go-user/handlers"
)

func main() {
	http.HandleFunc("/users", handlers.HandleUsers)
	http.HandleFunc("/users/", handlers.HandleUserByID)

	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
