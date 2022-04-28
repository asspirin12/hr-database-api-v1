package handlers

import (
	"encoding/json"
	"hr-database-api/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Employees struct{}

func (e Employees) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/employees/" {
		if r.Method == http.MethodGet {
			e.getEmployees(rw, r)
			return
		} else if r.Method == http.MethodPost {
			e.addEmployee(rw, r)
			return
		}
	} else {
		path := strings.Trim(r.URL.Path, "/")
		pathElems := strings.Split(path, "/")
		if len(pathElems) < 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(pathElems[1])
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			e.getEmployeeById(rw, r, id)
		} else if r.Method == http.MethodDelete {
			// TODO Implement me
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

func (e Employees) getEmployeeById(rw http.ResponseWriter, r *http.Request, id int) {
	log.Println("Handle GET one employee")

	employee, err := models.GetEmployeeById(id)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(employee)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func (e Employees) getEmployees(rw http.ResponseWriter, r *http.Request) {
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

func (e Employees) addEmployee(rw http.ResponseWriter, r *http.Request) {
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
