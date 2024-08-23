package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connection user=postgres password=LALA!kaswag1234")
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect to : %v\n", err))

	}
	defer conn.Close()
	log.Println("connected to db")

	//test my connection

	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot ping database")
	}
	log.Println("pinged database")
	//get rows from tables
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// insert a row

	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(query, "Jack", "brown")
	if err != nil {
		log.Fatal(err)

	}
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	///udpate row

	stmt := `update users set first_name = $1 where first_name = $2`
	_, err = conn.Exec(stmt, "jackie", "jack")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("updated more rows")

	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	//get row by id

	query = `select first_name , last_name from users where first_name = $1`
	var first_name, last_name string
	row := conn.QueryRow(query, "Jane")
	err = row.Scan(&first_name, &last_name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("query row returns", first_name, last_name)

	// delete a row

	query = `delete from users where first_name = $1`
	_, err = conn.Exec(query, "Jack")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("delete a row ")
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("SELECT first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()
	var first_name, last_name string
	//var id int

	for rows.Next() {
		err := rows.Scan(&first_name, &last_name)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("record is", first_name, last_name)
	}
	fmt.Println("record is", first_name, last_name)
	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows")
	}
	fmt.Println("-")
	return nil
}
