package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
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
		return err
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

func GetEmployeeById(id int, rw http.ResponseWriter) (Employee, error) {
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
		return person, err
	}

	return person, nil
}

func GetEmployeesByDepartment(dep string) ([]Employee, error) {

	var department string
	depURIList := "\n/marketing\n/training\n/research_and_development\n/sales\n/business_development\n/product_management\n/support\n/legal\n/accounting\n/services\n/hr\n/engineering"

	switch dep {
	case "marketing":
		department = "Marketing"
	case "training":
		department = "Training"
	case "research_and_development":
		department = "Research and Development"
	case "sales":
		department = "Sales"
	case "business_development":
		department = "Business Development"
	case "product_management":
		department = "Product Management"
	case "support":
		department = "Support"
	case "legal":
		department = "Legal"
	case "accounting":
		department = "Accounting"
	case "services":
		department = "Services"
	case "hr":
		department = "Human Resources"
	case "engineering":
		department = "Engineering"
	default:
		return []Employee{}, fmt.Errorf("department \"%s\" not found, check department URI: %s", dep, depURIList)
	}

	statement := `
SELECT id, first_name, last_name, email, department, date_hired 
FROM employees 
WHERE department = $1`

	rows, err := DB.Query(statement, department)
	if err != nil {
		log.Fatal(err)
	}

	workforce := make([]Employee, 0)

	for rows.Next() {
		person := Employee{}
		rows.Scan(
			&person.Id,
			&person.FirstName,
			&person.LastName,
			&person.Email,
			&person.Department,
			&person.DateHired,
		)

		workforce = append(workforce, person)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return workforce, nil
}

func AddEmployee(newEmployee Employee) (int, error) {
	employeeId := 0

	statement := `
INSERT INTO employees (first_name, last_name, email, department, date_hired) 
VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := DB.QueryRow(statement,
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

func DeleteEmployee(id int) error {

	statement := `DELETE FROM employees WHERE id = $1`

	_, err := DB.Query(statement, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEmployee(employee Employee, id int) error {

	statement := `
UPDATE employees SET first_name = $1, last_name = $2, email = $3, department = $4, date_hired = $5
WHERE id = $6`

	_, err := DB.Query(statement,
		employee.FirstName,
		employee.LastName,
		employee.Email,
		employee.Department,
		employee.DateHired,
		id)

	return err
}
