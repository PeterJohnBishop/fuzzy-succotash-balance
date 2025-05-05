package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"fuzzy-succotash-balance/main.go/database"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, port string, db *sql.DB) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("Drinking Gin on %s", port),
		})
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
	r.POST("/drop/:table", func(c *gin.Context) {
		table := c.Param("table")
		database.DropTable(db, c, table)
	})
}

func addUserRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", func(c *gin.Context) {
		database.Login(db, c)
	})
	r.POST("/register", func(c *gin.Context) {
		database.CreateUser(db, c)
	})
	r.GET("/users", func(c *gin.Context) {
		database.GetUsers(db, c)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		database.GetUserByID(db, c)
	})
	r.PUT("/users/:id", func(c *gin.Context) {
		database.UpdateUserByID(db, c)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		database.DeleteUserByID(db, c)
	})
}

func addProductRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/products", func(c *gin.Context) {
		database.GetProducts(db, c)
	})
	r.POST("/products", func(c *gin.Context) {
		database.CreateProduct(db, c)
	})
	r.GET("/products/:upc", func(c *gin.Context) {
		database.GetProductByUPC(db, c)
	})
	r.PUT("/products/:upc", func(c *gin.Context) {
		database.UpdateProductByUPC(db, c)
	})
	r.DELETE("/products/:upc", func(c *gin.Context) {
		database.DeleteProductByUPC(db, c)
	})
}

func addOrderRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/orders", func(c *gin.Context) {
		database.GetOrders(db, c)
	})
	r.POST("/orders", func(c *gin.Context) {
		database.CreateOrder(db, c)
	})
	r.GET("/orders/:orderNumber", func(c *gin.Context) {
		database.GetOrderByNumber(db, c)
	})
	r.PUT("/orders/:orderNumber", func(c *gin.Context) {
		database.UpdateOrderByNumber(db, c)
	})
	r.DELETE("/orders/:orderNumber", func(c *gin.Context) {
		database.DeleteOrderByNumber(db, c)
	})
}

func addChatMessageingRoutes(r *gin.Engine, db *sql.DB) {

	r.POST("/chats", func(c *gin.Context) {
		database.CreateChat(db, c)
	})
	r.POST("/messages", func(c *gin.Context) {
		database.CreateMessage(db, c)
	})
	r.GET("/chats", func(c *gin.Context) {
		database.GetAllChats(db, c)
	})
	r.GET("/chats/:chatID", func(c *gin.Context) {
		database.GetChatByID(db, c)
	})
	r.GET("/chats/:chatID/messages", func(c *gin.Context) {
		database.GetChatWithMessages(db, c)
	})
	r.PUT("/chats/:chatID", func(c *gin.Context) {
		database.UpdateChatByID(db, c)
	})

	r.DELETE("/chats/:chatID", func(c *gin.Context) {
		database.DeleteChatByID(db, c)
	})
	r.DELETE("/messages/:messageID", func(c *gin.Context) {
		database.DeleteMessageByID(db, c)
	})
}
