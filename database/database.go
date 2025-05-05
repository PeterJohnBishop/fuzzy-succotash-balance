package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
	fmt.Println("Connecting with:", psqlInfo)

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

func DropTable(db *sql.DB, c *gin.Context, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, tableName)

	_, err := db.Exec(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	log.Printf("Table %s dropped successfully.\n", tableName)
	c.JSON(http.StatusOK, gin.H{"message": "Table dropped successfully"})
	return nil
}

func CreateTable(db *sql.DB, c *gin.Context, tableName string) {
	switch tableName {
	case "products":
		err := CreateProductsTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "orders":
		err := CreateOrdersTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "users":
		err := CreateUsersTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "chats":
		err := CreateChatsTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "messages":
		err := CreateMessagesTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table name"})
		return
	}
}
