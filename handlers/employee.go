package handlers

import (
	"encoding/json"
	"hr-database-api/models"
	"log"
	"net/http"
)

type Employees struct{}

func (e Employees) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		e.GetEmployees(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		e.AddEmployee(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (e Employees) GetEmployees(rw http.ResponseWriter, r *http.Request) {
	log.Println("Handle GET employees")

	employeesList, err := models.GetEmployees(10)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(rw)

	err = encoder.Encode(employeesList)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func (e Employees) AddEmployee(rw http.ResponseWriter, r *http.Request) {
	log.Println("Handle POST employee")

	employee := models.Employee{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&employee)
	if err != nil {
		http.Error(rw, "Unable to decode json", http.StatusInternalServerError)
		log.Fatal(err)
	}

	id, err := models.AddEmployee(employee)
	if err != nil {
		http.Error(rw, "Failed to add a new employee", http.StatusInternalServerError)
		log.Fatal(err)
	}
	log.Printf("Added a record with the id %d", id)
}
