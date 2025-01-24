package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IDOMATH/docker-example/db"
	"github.com/IDOMATH/docker-example/util"
)

var DS db.DataStore

func main() {
	serverPort, err := util.EnvMust("CONTAINER_PORT") //":8080"
	if err != nil {
		panic(err.Error())
	}
	dbPort, err := util.EnvMust("DB_PORT") //"5432"
	if err != nil {
		panic(err.Error())
	}
	dbname, err := util.EnvMust("POSTGRES_DB") //"docker-example"
	if err != nil {
		panic(err.Error())
	}
	user, err := util.EnvMust("POSTGRES_USER") // "postgres"
	if err != nil {
		panic(err.Error())
	}
	password, err := util.EnvMust("POSTGRES_PASSWORD") // "mysecretpassword"
	if err != nil {
		panic(err.Error())
	}
	sslmode := "disable"

	router := http.NewServeMux()

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /", handlePostData)

	router.HandleFunc("GET /{id}", handleGetById)
	router.HandleFunc("PUT /{id}", handlePutData)
	router.HandleFunc("DELETE /{id}", handleDeleteData)
	router.HandleFunc("GET /data", handleGetAllData)

	router.HandleFunc("POST /seed", handleInit)
	router.HandleFunc("POST /drop", handleDrop)
	router.HandleFunc("POST /reset", handleResetDb)

	server := http.Server{
		Addr:    ":" + serverPort,
		Handler: router,
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@db:%s/%s?sslmode=%s", user, password, dbPort, dbname, sslmode)

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

func handleInit(w http.ResponseWriter, r *http.Request) {
	err := DS.InitDb()
	if err != nil {
		w.Write([]byte("error seeding db"))
		return
	}
	w.Write([]byte("Seeding db"))
}

func handleDrop(w http.ResponseWriter, r *http.Request) {
	err := DS.Drop()
	if err != nil {
		w.Write([]byte("eror dropping db"))
		return
	}
	w.Write([]byte("Dropping db"))
}

func handleResetDb(w http.ResponseWriter, r *http.Request) {
	err := DS.Drop()
	if err != nil {
		w.Write([]byte("eror dropping db"))
		return
	}
	w.Write([]byte("Dropping db"))

	err = DS.InitDb()
	if err != nil {
		w.Write([]byte("error seeding db"))
		return
	}
	w.Write([]byte("Seeding db"))
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
		w.Write([]byte(datum.Data + "\n"))
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
