package main

import (
	"log"
	"net/http"

	"hr-database-api/data"
	"hr-database-api/handlers"

	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := data.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	e := handlers.Employees{}

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
