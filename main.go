package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IDOMATH/docker-example/db"
)

var DS db.DataStore

func main() {
	fmt.Println("Hello world")
	router := http.NewServeMux()

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /", handlePostData)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", "localhost", "5432", "postgres", "postgres", "mysecretpassword", "disable")
	postgresDb, err := db.ConnectSql(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	DS = *db.NewDataStore(postgresDb.Sql)

	fmt.Println("Starting on port 8080")
	log.Fatal(server.ListenAndServe())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home"))
}

func handlePostData(w http.ResponseWriter, r *http.Request) {
	DS.InsertData(time.Now().String())
	w.Write([]byte("Entered"))
}
