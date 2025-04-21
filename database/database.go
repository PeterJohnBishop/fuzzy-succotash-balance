package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func ConnectPSQL(db *sql.DB) *sql.DB {

	host := os.Getenv("PSQL_HOST")
	portStr := os.Getenv("PSQL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid port number:", err)
	}
	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASSWORD")
	dbname := os.Getenv("PSQL_DBNAME")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	mydb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = mydb.Ping()
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to Postgres container on :%s", portStr)
	return mydb
}
