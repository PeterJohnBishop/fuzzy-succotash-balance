package server

import (
	"fmt"
	"log"
	"net/http"
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

	log.Printf("Listening on %s", port)
	r.Run(port)
}

func setupRoutes(r *gin.Engine, port string) {

	r.GET("/", func(c *gin.Context) {
		response := map[string]interface{}{
			"msg": fmt.Sprintf("Drinking Gin on %s", port),
		}
		c.JSON(http.StatusOK, response)
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	r.GET("/apple-touch-icon.png", func(c *gin.Context) {
		c.Status(204)
	})

	r.GET("/apple-touch-icon-precomposed.png", func(c *gin.Context) {
		c.Status(204)
	})

}
