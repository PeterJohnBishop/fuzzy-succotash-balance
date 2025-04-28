package database

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID NOT NULL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		avatar TEXT,
		online BOOL DEFAULT false,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	return err
}

func SeedNewUsers(count int, db *sql.DB) error {
	err := CreateUsersTable(db)
	if err != nil {
		return err
	}
	err = SetFakeUsers(count, db)
	return err
}

func CreateUser(db *sql.DB, c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO users (id, name, email, password, avatar, online, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`
	_, err := db.ExecContext(c, query, user.ID, user.Name, user.Email, user.Password, user.Avatar, user.Online)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}

func GetUsers(db *sql.DB, c *gin.Context) {
	rows, err := db.QueryContext(c, "SELECT id, name, email, password, avatar, online, created_at, updated_at FROM users;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.Online, &user.CreatedAt, &user.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var user User
	query := `SELECT id, name, email, password, avatar, online, created_at, updated_at FROM users WHERE id = $1`
	err := db.QueryRowContext(c, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.Online, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE users SET name=$1, email=$2, password=$3, avatar=$4, online=$5, updated_at=NOW() WHERE id=$6`
	result, err := db.ExecContext(c, query, user.Name, user.Email, user.Password, user.Avatar, user.Online, id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated!"})
}

func DeleteUserByID(db *sql.DB, c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.ExecContext(c, query, id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted!"})
}
