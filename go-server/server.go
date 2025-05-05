package server

import (
	"database/sql"
	"log"
	"os"

	"fuzzy-succotash-balance/main.go/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func StartServer(db *sql.DB) {
	log.Println("Starting Server container")

	port := os.Getenv("GIN_PORT")

	r := gin.Default()
	r.Use(database.VerifyJWT())
	err := r.SetTrustedProxies([]string{"172.16.0.0/12"})
	if err != nil {
		log.Fatal(err)
	}

	setupRoutes(r, port, db)
	addUserRoutes(r, db)
	addProductRoutes(r, db)
	addOrderRoutes(r, db)
	addChatMessageingRoutes(r, db)

	r.Run(port)
}
