package database

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Avatar    string    `json:"avatar"`
	Online    bool      `json:"online"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	UPC         string   `json:"upc"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type OrderStatus string

const (
	NotSent    OrderStatus = "Not Sent"
	Sent       OrderStatus = "Sent"
	Received   OrderStatus = "Received"
	InProgress OrderStatus = "In Progress"
	InTransit  OrderStatus = "In Transit"
	Delivered  OrderStatus = "Delivered"
)

type Order struct {
	OrderNumber int         `json:"orderNumber"`
	Status      OrderStatus `json:"status"`
	User        string      `json:"user"`
	Products    []string    `json:"products"`
	Total       float64     `json:"total"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

type Chat struct {
	ChatID    string    `json:"chat_id"`
	Users     []string  `json:"users"`
	Messages  []string  `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	MessageID string    `json:"message_id"`
	Chat      string    `json:"chat"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Media     []string  `json:"media"`
	CreatedAt time.Time `json:"created_at"`
}
