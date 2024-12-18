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
	serverPort := ":8080"

	// host := "localhost"
	dbPort := "5432"
	dbname := "docker-example"
	user := "postgres"
	password := "mysecretpassword"
	sslmode := "disable"

	fmt.Println("Hello world")
	router := http.NewServeMux()

	router.HandleFunc("GET /", setHeaders(handleHome))
	router.HandleFunc("POST /", handlePostData)
	router.HandleFunc("OPTIONS /", handlePreflight)
	router.HandleFunc("GET /{id}", handleGetById)
	router.HandleFunc("PUT /{id}", handlePutData)
	router.HandleFunc("DELETE /{id}", handleDeleteData)

	server := http.Server{
		Addr:    serverPort,
		Handler: router,
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=%s", user, password, dbPort, dbname, sslmode)
	// connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, dbPort, dbname, user, password, sslmode)
	postgresDb, err := db.ConnectSql(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to postgres on port: ", serverPort)
	DS = *db.NewDataStore(postgresDb.Sql)

	fmt.Println("Starting on port ", serverPort)
	log.Fatal(server.ListenAndServe())
}

func setHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("headers set"))
		next(w, r)
	}
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	fmt.Println([]byte("options"))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	w.Write([]byte("Welcome Home"))
}

func handlePostData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("posting...")
	DS.InsertData(time.Now().String())
	w.Write([]byte("Entered"))
}

func handleGetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte("error parsing id"))
		return
	}
	fmt.Println("getting by id")
	entry, err := DS.GetDataById(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(entry.Data))
}

func handlePutData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error parsing id to int: %d", id)))
		return
	}
	err = DS.UpdateData(db.Entry{Id: id, Data: time.Now().String()})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("updated"))
}

func handleDeleteData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error parsing id to int: %d", id)))
		return
	}
	err = DS.DeleteData(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("deleted"))
}
