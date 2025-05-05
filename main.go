package main

import (
	"database/sql"

	"fuzzy-succotash-balance/main.go/database"
	"fuzzy-succotash-balance/main.go/server"
	"log"
)

var db *sql.DB

func main() {
	log.Println("Starting Fuzzy-Succotash-Balance")
	db := database.ConnectPSQL(db)
	database.SeedNewUsers(500, db)
	database.CreateUpdatedAtTrigger(db)
	database.CreateUpdatedAtTriggerForTable(db, "users")
	server.StartServer(db)
}
