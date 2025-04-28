package main

import (
	"database/sql"
	"fmt"

	"fuzzy-succotash-balance/main.go/database"
	"fuzzy-succotash-balance/main.go/server"
	"log"
)

var db *sql.DB

func main() {
	log.Println("Starting Fuzzy-Succotash-Balance")
	db := database.ConnectPSQL(db)
	// DropTable(db, "users")
	database.SeedNewUsers(500, db)
	database.CreateUpdatedAtTrigger(db)
	database.CreateUpdatedAtTriggerForTable(db, "users")
	server.StartServer(db)
}

func DropTable(db *sql.DB, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, tableName)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to drop table %s: %w", tableName, err)
	}

	log.Printf("Table %s dropped successfully.\n", tableName)
	return nil
}
