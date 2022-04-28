package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 54333
	user     = "guest"
	password = "guest"
	dbname   = "guest"
)

func ConnectDatabase() error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to database")
	DB = db
	return nil
}

type Employee struct {
	Id         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Department string `json:"department"`
	Email      string `json:"email,omitempty"`
	DateHired  string `json:"date_hired"`
}

func GetEmployeeById(id int) (Employee, error) {
	statement := `SELECT first_name, last_name, email, department, date_hired FROM employees WHERE id = $1`
	row := DB.QueryRow(statement, id)

	person := Employee{Id: id}

	err := row.Scan(
		&person.FirstName,
		&person.LastName,
		&person.Email,
		&person.Department,
		&person.DateHired,
	)

	if err != nil {
		log.Fatal(err)
	}

	return person, nil
}

func GetEmployees(count int) ([]Employee, error) {
	statement := `SELECT id, first_name, last_name, email, department, date_hired FROM employees LIMIT ` + strconv.Itoa(count)

	rows, err := DB.Query(statement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	workforce := make([]Employee, 0)

	for rows.Next() {
		person := Employee{}
		err = rows.Scan(
			&person.Id,
			&person.FirstName,
			&person.LastName,
			&person.Email,
			&person.Department,
			&person.DateHired,
		)
		if err != nil {
			log.Fatal(err)
		}

		workforce = append(workforce, person)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return workforce, err
}

func AddEmployee(newEmployee Employee) (int, error) {
	employeeId := 0
	err := DB.QueryRow(`
INSERT INTO employees (first_name, last_name, email, department, date_hired) 
VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		newEmployee.FirstName,
		newEmployee.LastName,
		newEmployee.Email,
		newEmployee.Department,
		newEmployee.DateHired,
	).Scan(&employeeId)

	if err != nil {
		return 0, err
	}

	return employeeId, nil
}
