package database

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateProductsTable(db *sql.DB) error {
	query := `
	CREATE TABLE products (
		upc TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price FLOAT,
		images TEXT[], 
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	return err
}

func CreateProduct(db *sql.DB, c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO products (upc, name, description, price, images, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`
	_, err := db.ExecContext(c, query, product.UPC, product.Name, product.Description, product.Price, pq.Array(product.Images))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created!"})
}

func GetProducts(db *sql.DB, c *gin.Context) {
	rows, err := db.QueryContext(c, "SELECT upc, name, description, price, images, created_at, updated_at FROM products;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		var images pq.StringArray

		if err := rows.Scan(&product.UPC, &product.Name, &product.Description, &product.Price, &images, &product.CreatedAt, &product.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		product.Images = []string(images) // Convert pq.StringArray to []string
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByUPC(db *sql.DB, c *gin.Context) {
	upc := c.Param("upc")
	var product Product
	query := `SELECT upc, name, description, price, created_at, updated_at FROM products WHERE upc = $1`
	err := db.QueryRowContext(c, query, upc).Scan(&product.UPC, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProductByUPC(db *sql.DB, c *gin.Context) {
	upc := c.Param("upc")
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE products SET name=$1, description=$2, price=$3, updated_at=NOW() WHERE upc=$4`
	result, err := db.ExecContext(c, query, product.Name, product.Description, product.Price, upc)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated!"})
}

func DeleteProductByUPC(db *sql.DB, c *gin.Context) {
	upc := c.Param("upc")
	query := `DELETE FROM products WHERE upc = $1`
	result, err := db.ExecContext(c, query, upc)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted!"})
}
