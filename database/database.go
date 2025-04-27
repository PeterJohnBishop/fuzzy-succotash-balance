package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectPSQL(db *sql.DB) *sql.DB {

	host := "postgres"
	port := 5432
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
	log.Printf("Connected to Postgres container on :%d", port)
	return mydb
}

func CreateUpdatedAtTrigger(db *sql.DB) error {
	triggerFunc := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	if _, err := db.Exec(triggerFunc); err != nil {
		return fmt.Errorf("error creating trigger function: %w", err)
	}

	return nil
}

func CreateUpdatedAtTriggerForTable(db *sql.DB, tableName string) error {
	trigger := fmt.Sprintf(`
	CREATE TRIGGER update_%s_updated_at
	BEFORE UPDATE ON %s
	FOR EACH ROW
	EXECUTE PROCEDURE update_updated_at_column();`, tableName, tableName)

	_, err := db.Exec(trigger)
	if err != nil {
		return fmt.Errorf("could not create trigger for table %s: %w", tableName, err)
	}
	return nil
}
