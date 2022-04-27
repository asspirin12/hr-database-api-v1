create table employees (
       id serial primary key,
       first_name varchar(50),
       last_name varchar(50),
       email varchar(50),
       department varchar(50),
       date_hired date
);