package database

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrdersTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS orders (
        order_number SERIAL PRIMARY KEY,
        status TEXT NOT NULL,
        user_id TEXT NOT NULL,
        products JSONB NOT NULL,
        total FLOAT8 NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`

	_, err := db.Exec(query)
	return err
}

func CreateOrder(db *sql.DB, c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productsJSON, err := json.Marshal(order.Products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO orders (status, user_id, products, total, created_at, updated_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING order_number`

	err = db.QueryRowContext(c, query, order.Status, order.User, productsJSON, order.Total).Scan(&order.OrderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created!", "orderNumber": order.OrderNumber})
}

func GetOrders(db *sql.DB, c *gin.Context) {
	query := `SELECT order_number, status, user_id, products, total, created_at, updated_at FROM orders`
	rows, err := db.QueryContext(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		var productsData []byte
		if err := rows.Scan(&order.OrderNumber, &order.Status, &order.User, &productsData, &order.Total, &order.CreatedAt, &order.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		json.Unmarshal(productsData, &order.Products)
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrderByNumber(db *sql.DB, c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	query := `SELECT order_number, status, user_id, products, total, created_at, updated_at FROM orders WHERE order_number = $1`

	var order Order
	var productsData []byte
	err := db.QueryRowContext(c, query, orderNumber).Scan(&order.OrderNumber, &order.Status, &order.User, &productsData, &order.Total, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	json.Unmarshal(productsData, &order.Products)
	c.JSON(http.StatusOK, order)
}

func UpdateOrderByNumber(db *sql.DB, c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productsJSON, err := json.Marshal(order.Products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE orders 
              SET status=$1, user_id=$2, products=$3, total=$4, updated_at=NOW()
              WHERE order_number=$5`

	result, err := db.ExecContext(c, query, order.Status, order.User, productsJSON, order.Total, orderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated!"})
}

func DeleteOrderByNumber(db *sql.DB, c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	query := `DELETE FROM orders WHERE order_number = $1`

	result, err := db.ExecContext(c, query, orderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted!"})
}
