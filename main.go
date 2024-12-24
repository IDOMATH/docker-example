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

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /", handlePostData)
	router.HandleFunc("GET /{id}", handleGetById)
	router.HandleFunc("PUT /{id}", handlePutData)
	router.HandleFunc("DELETE /{id}", handleDeleteData)
	router.HandleFunc("POST /seed", handleSeed)

	router.HandleFunc("GET /data", handleGetAllData)

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

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home"))
}

func handlePostData(w http.ResponseWriter, r *http.Request) {
	DS.InsertData(time.Now().String())
	w.Write([]byte("Entered"))
}

func handleGetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Write([]byte("error parsing id"))
		return
	}
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

func handleGetAllData(w http.ResponseWriter, r *http.Request) {
	data, err := DS.GetAllData()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	for _, datum := range data {
		w.Write([]byte(datum.Data))
	}
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
