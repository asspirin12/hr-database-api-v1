package main

import (
	_ "github.com/lib/pq"
	"hr-database-api/handlers"
	"hr-database-api/models"
	"log"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := models.ConnectDatabase()
	checkError(err)

	e := handlers.Employees{}
	mux := http.NewServeMux()
	mux.Handle("/employees", e)

	s := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err = s.ListenAndServe()
	checkError(err)
}
