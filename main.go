package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello world")
	router := http.NewServeMux()

	router.HandleFunc("GET /", handleHome)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting on port 8080")
	log.Fatal(server.ListenAndServe())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home"))
}
