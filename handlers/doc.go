// Package handlers. Package classification of HR Database API
//
// Documentation for HR Database API
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
//  - application/json
//
// Produces:
//  - application/json
// swagger:meta
package handlers

import "hr-database-api/data"

// An error message returned as a string and an HTTP status code
// swagger:response errorResponse
type errorResponse struct {
	// Description of the error
	// in: body
	Body Error
}

// A list of employees returned
// swagger:response employeesResponse
type employeesResponse struct {
	// All employees in the database
	// in: body
	Body []data.Employee
}

// A single employee returned
// swagger:response employeeResponse
type employeeResponse struct {
	// An employee
	// in: body
	Body data.Employee
}

//
// swagger:parameters addEmployee updateEmployee
type employeeParams struct {
	// Employee data structure to Add or Update
	// The id field is ignored by Add operation
	// in: body
	// required: true
	Body data.Employee
}

// swagger:parameters updateEmployee deleteEmployee getEmployeeById
type employeeId struct {
	// The id of the employee whose record is retrieved, updated or deleted
	// in: path
	// required: true
	Id int `json:"id"`
}

// swagger:parameters getEmployeesByDepartment
type department struct {
	// The department name
	// in: path
	// required: true
	Department string `json:"department"`
}
