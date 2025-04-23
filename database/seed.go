package database

import (
	"database/sql"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func SetFakeUsers(count int, db *sql.DB) error {
	gofakeit.Seed(0)

	for i := 0; i < count; i++ {

		var newUser User = User{
			ID:        gofakeit.UUID(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, true, false, 16),
			Avatar:    gofakeit.ImageURL(500, 500),
			Online:    false,
			CreatedAt: time.Now(),
		}

		hashedPassword, err := HashedPassword(newUser.Password)
		if err != nil {
			return err
		}

		query := "INSERT INTO users (id, name, email, password, avatar, online) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at"
		err = db.QueryRow(query, newUser.ID, newUser.Name, newUser.Email, hashedPassword, newUser.Avatar, newUser.Online).Scan(&newUser.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}
