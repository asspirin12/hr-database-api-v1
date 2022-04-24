package handlers

import (
	"encoding/json"
	"hr-database-api/models"
	"log"
	"net/http"
)

type Employees struct{}

func (e Employees) GetEmployees(rw http.ResponseWriter, r *http.Request) {
	log.Println("Handle GET employees")

	employeesList, err := models.GetEmployees(10)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := json.Marshal(employeesList)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	_, err = rw.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func (e Employees) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		e.GetEmployees(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}
