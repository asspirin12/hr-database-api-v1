package handlers

import (
	"encoding/json"
	"hr-database-api/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Error error
// swagger:model
type Error struct {
	// message
	Message string `json:"message"`
	// code
	Code int `json:"code"`
}

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

// swagger:route GET /employees/ employees listEmployees
// Return a list of employees
// responses:
//  200: employeesResponse
//  404: errorResponse
//  405: errorResponse

// getEmployees returns a list of employees (limit is 10 records by default)
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

// swagger:route GET /employees/{id} employees getEmployeeById
// Return a single employee by id
// responses:
//  200: employeeResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse

// getEmployeeById returns a single employee which matches the id
func (e Employees) getEmployeeById(rw http.ResponseWriter, r *http.Request, id int) {
	log.Println("Handle GET one employee")

	employee, err := models.GetEmployeeById(id, rw)
	if err != nil {
		http.Error(rw, "Unable to retrieve data, check the id", http.StatusNotFound)
		log.Println(err)
	}

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(employee)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
		log.Println(err)
	}
}

// swagger:route GET /employees/{department} employees getEmployeesByDepartment
// Return a list of employees working in a particular department
// responses:
//  200: employeesResponse
//  404: errorResponse
//  405: errorResponse

// getEmployeesByDepartment returns a list of employees working in a particular department
func (e Employees) getEmployeesByDepartment(rw http.ResponseWriter, r *http.Request, dep string) {
	log.Println("Handle GET employees by department")
	employeesByDepartment, err := models.GetEmployeesByDepartment(dep)
	if err != nil {
		http.Error(rw, "Department not found, check URI", http.StatusNotFound)
		log.Println(err)
	}

	encoder := json.NewEncoder(rw)

	err = encoder.Encode(employeesByDepartment)
	if err != nil {
		http.Error(rw, "Unable to encode json", http.StatusInternalServerError)
		log.Println(err)
	}
}

// swagger:route POST /employees/ employees addEmployee
// Add a new employee to the end of the database
// responses:
//  500: errorResponse

// addEmployee adds an employee to the end of the database
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

// swagger:route POST /employees/{id} employees updateEmployee
// Update an employee record
// responses:
//  500: errorResponse

// updateEmployee updates an employee record
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
		http.Error(rw, "Failed to update a record", http.StatusInternalServerError)
		log.Println(err)
	}
}

// swagger:route DELETE /employees/{id} employees deleteEmployee
// Delete an employee from the database
// responses:
//  500: errorResponse

// deleteEmployee deletes an employee from the database by id
func (e Employees) deleteEmployee(rw http.ResponseWriter, r *http.Request, id int) {
	log.Println("Handle DELETE employee")

	err := models.DeleteEmployee(id)
	if err != nil {
		http.Error(rw, "Failed to delete a record, check id", http.StatusNotFound)
		log.Println(err)
	}
}
