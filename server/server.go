package server

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func StartServer() {
	log.Println("Starting Server")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("GIN_PORT")

	r := gin.Default()
	setupRoutes(r, port)
	loadTestingRoutes(r)

	log.Printf("Listening on %s", port)
	r.Run(port)
}
