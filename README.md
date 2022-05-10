# HR Database Demo

The application presents the following REST API to clients:

`GET /employees/` returns a list of employees (by default the limit is 10 records)\
`GET /employees/<id>` returns a single record by id\
`POST /employees/` adds an employee to the end of the database\
`POST /employees/<id>` updates existing record\
`GET /department/<department>` returns a list of employees by department. The list of URIs for departments: marketing, training, research_and_development, sales, business_development, product_management, support, legal, accounting, services, hr, engineering.




## How to work with the project in GoLand

### Connect to database 

Start a container in `docker-compose.yml,` connect to the container from the Database tool window (copy a driver comment entry `jdbc:postgresql://localhost:54333/guest?user=guest&password=guest` from `docker-compose.yml` and paste it to URL field in Data Sources and Drivers pop-up.). 

Create the database structure by running `storage/init.sql`. Populate the database with the fake data by running `storage/employees.sql`. 

### Start the server 

Run `main.go`. 

### Run the requests from the `http-request.http`

