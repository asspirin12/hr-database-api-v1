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
		} else if r.Method == http.MethodPost {
			e.addEmployee(rw, r)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			log.Printf("Method %s is not allowed for %s", r.Method, r.URL.Path)
			return
		}
		// if URL starts with /department/ followed by a department name
	} else if strings.HasPrefix(r.URL.Path, "/department/") {
		if r.Method == http.MethodGet {
			dep := extractPathElem(rw, r)
			e.getEmployeesByDepartment(rw, r, dep)
			return
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			log.Printf("Method %s is not allowed for %s", r.Method, r.URL.Path)
			return
		}
		// if URL starts with /employees/ followed by id number
	} else if strings.HasPrefix(r.URL.Path, "/employees/") {
		// extract id from URL
		idString := extractPathElem(rw, r)

		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Unable to parse id", http.StatusBadRequest)
			log.Println("Unable to parse id, error: ", err)
			return
		}

		if r.Method == http.MethodGet {
			e.getEmployeeById(rw, r, id)
		} else if r.Method == http.MethodDelete {
			e.deleteEmployee(rw, r, id)
		} else if r.Method == http.MethodPost {
			e.updateEmployee(rw, r, id)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			log.Printf("Method %s is not allowed for %s", r.Method, r.URL.Path)
			return
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
		log.Fatalf("Page %s is not found", r.URL.Path)
		return
	}
}

// extracts id from URL
func extractPathElem(rw http.ResponseWriter, r *http.Request) string {
	path := strings.Trim(r.URL.Path, "/")
	pathElems := strings.Split(path, "/")
	if len(pathElems) < 2 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return ""
	}

	return pathElems[1]
}

func (e Employees) getEmployeeById(rw http.ResponseWriter, r *http.Request, id int) {
	log.Println("Handle GET one employee")

	employee, err := models.GetEmployeeById(id, rw)
	if err != nil {
		http.Error(rw, "Unable to retrieve data, check the id", http.StatusBadRequest)
		log.Println(err)
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(employee)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusBadRequest)
		log.Println(err)
	}
}

func (e Employees) getEmployeesByDepartment(rw http.ResponseWriter, r *http.Request, dep string) {
	log.Println("Handle GET employees by department")
	employeesByDepartment, err := models.GetEmployeesByDepartment(dep)
	if err != nil {
		http.Error(rw, "Department not found, check URI", http.StatusBadRequest)
		log.Println(err)
	}

	encoder := json.NewEncoder(rw)

	err = encoder.Encode(employeesByDepartment)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (e Employees) getEmployees(rw http.ResponseWriter, r *http.Request) {
	log.Println("Handle GET employees")

	employeesList, err := models.GetEmployees(10)
	if err != nil {
		http.Error(rw, "Unable to get employees list", http.StatusInternalServerError)
		log.Println(err)
	}

	encoder := json.NewEncoder(rw)

	err = encoder.Encode(employeesList)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (e Employees) addEmployee(rw http.ResponseWriter, r *http.Request) {
	log.Println("Handle POST employee")

	employee := models.Employee{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&employee)
	if err != nil {
		http.Error(rw, "Unable to decode json", http.StatusInternalServerError)
		log.Println(err)
	}

	id, err := models.AddEmployee(employee)
	if err != nil {
		http.Error(rw, "Failed to add a new employee", http.StatusInternalServerError)
		log.Println(err)
	}
	log.Printf("Added a record with the id %d", id)
}

func (e Employees) deleteEmployee(rw http.ResponseWriter, r *http.Request, id int) {
	log.Println("Handle DELETE employee")

	err := models.DeleteEmployee(id)
	if err != nil {
		http.Error(rw, "Failed to delete a record", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (e Employees) updateEmployee(rw http.ResponseWriter, r *http.Request, id int) {
	log.Println("Handle UPDATE employee")

	employee := models.Employee{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&employee)
	if err != nil {
		http.Error(rw, "Unable to decode json", http.StatusInternalServerError)
		log.Println(err)
	}

	err = models.UpdateEmployee(employee, id)
	if err != nil {
		http.Error(rw, "Failed to update a record", http.StatusBadRequest)
		log.Println(err)
	}
}
