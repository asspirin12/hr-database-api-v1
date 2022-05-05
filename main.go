package main

import (
	_ "github.com/lib/pq"
	"hr-database-api/handlers"
	"hr-database-api/models"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := models.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	e := handlers.Employees{}
	mux := http.NewServeMux()

	mux.Handle("/employees/", e)
	mux.Handle("/department/", e)

	s := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
