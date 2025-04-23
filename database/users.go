package database

import (
	"database/sql"

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

func GetUsers(db *sql.DB, c *gin.Context) ([]User, error) {

	rows, err := db.QueryContext(c, "SELECT id, name, email, password, avatar, online, created_at, updated_at FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Avatar, &user.Online, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return []User{}, err
	}
	return users, nil
}
