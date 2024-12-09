package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IDOMATH/docker-example/db"
)

var DS db.DataStore

func main() {
	fmt.Println("Hello world")
	router := http.NewServeMux()

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /", handlePostData)
	router.HandleFunc("PUT /{id}", handlePutData)
	router.HandleFunc("DELETE /{id}", handleDeleteData)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", "localhost", "5432", "docker-example", "postgres", "mysecretpassword", "disable")
	postgresDb, err := db.ConnectSql(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	DS = *db.NewDataStore(postgresDb.Sql)

	fmt.Println("Starting on port 8080")
	log.Fatal(server.ListenAndServe())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	w.Write([]byte("Welcome Home"))
}

func handlePostData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("posting...")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	DS.InsertData(time.Now().String())
	w.Write([]byte("Entered"))
}

func handlePutData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error parsing id to int: %d", id)))
	}
	err = DS.UpdateData(db.Entry{Id: id, Data: time.Now().String()})
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("updated"))
}

func handleDeleteData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error parsing id to int: %d", id)))
	}
	err = DS.DeleteData(id)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("deleted"))
}
