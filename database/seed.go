package database

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
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

func SetFakeProducts(count int, db *sql.DB) error {
	gofakeit.Seed(0)
	existingUPCs := make(map[string]bool)

	for i := 0; i < count; i++ {
		var newProduct Product = Product{
			UPC:         generateFakeUPC(existingUPCs),
			Name:        gofakeit.ProductName(),
			Description: gofakeit.Paragraph(1, 3, 5, " "),
			Price:       gofakeit.Price(10, 1000), // Random price between 10 and 1000
			Images: []string{
				gofakeit.ImageURL(600, 600),
				gofakeit.ImageURL(600, 600),
				gofakeit.ImageURL(600, 600),
			},
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		query := `INSERT INTO products (upc, name, description, price, images, created_at, updated_at)
				  VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := db.Exec(
			query,
			newProduct.UPC,
			newProduct.Name,
			newProduct.Description,
			newProduct.Price,
			pq.Array(newProduct.Images),
			newProduct.CreatedAt,
			newProduct.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateFakeUPC(existing map[string]bool) string {
	for {
		number := rand.Int63n(899999999999) + 100000000000
		upc := strconv.FormatInt(number, 10)

		if !existing[upc] {
			existing[upc] = true
			return upc
		}
	}
}
