package database

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateChatsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS chats (
		chat_id TEXT PRIMARY KEY,
		users TEXT[] NOT NULL,
		messages TEXT[] DEFAULT '{}',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create chats table: %w", err)
	}
	return nil
}

func CreateMessagesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		message_id TEXT PRIMARY KEY,
		chat_id TEXT NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
		sender TEXT NOT NULL,
		text TEXT,
		media TEXT[] DEFAULT '{}',
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create messages table: %w", err)
	}
	return nil
}

func GenerateChatID(t time.Time) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const suffixLength = 6

	suffix := make([]byte, suffixLength)
	for i := range suffix {
		suffix[i] = charset[rand.Intn(len(charset))]
	}

	// Format: c_YYYYMMDD_HHMMSS_random
	return fmt.Sprintf("c_%s_%s",
		t.Format("20060102_150405"),
		string(suffix),
	)
}

func GenerateMessageID(senderUUID string) (string, error) {
	cleanUUID := strings.ReplaceAll(senderUUID, "-", "")

	// Take the first 8 characters
	if len(cleanUUID) < 8 {
		return "", fmt.Errorf("invalid sender UUID")
	}
	prefix := cleanUUID[:8]

	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomPart := hex.EncodeToString(randomBytes)

	timestamp := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS

	// Format: m_{senderprefix}_{timestamp}_{random}
	messageID := fmt.Sprintf("m_%s_%s_%s", prefix, timestamp, randomPart)

	return messageID, nil
}

func CreateChat(db *sql.DB, c *gin.Context) {
	var chat Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var time = time.Now()
	chat.ChatID = GenerateChatID(time)

	query := `INSERT INTO chats (chat_id, users, messages, created_at, updated_at)
	          VALUES ($1, $2, $3, NOW(), NOW())`

	_, err := db.ExecContext(c, query, chat.ChatID, pq.Array(chat.Users), pq.Array(chat.Messages))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Chat created successfully"})
}

func CreateMessage(db *sql.DB, c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageID, err := GenerateMessageID(msg.Sender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg.MessageID = messageID

	query := `INSERT INTO messages (message_id, chat_id, sender, text, media, created_at)
	          VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err = db.ExecContext(c, query, msg.MessageID, msg.Chat, msg.Sender, msg.Text, pq.Array(msg.Media))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message created successfully"})
}

func GetAllChats(db *sql.DB, c *gin.Context) {
	query := `SELECT chat_id, users, messages, created_at, updated_at FROM chats`
	rows, err := db.QueryContext(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ChatID, pq.Array(&chat.Users), pq.Array(&chat.Messages), &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		chats = append(chats, chat)
	}

	c.JSON(http.StatusOK, chats)
}

func GetChatWithMessages(db *sql.DB, c *gin.Context) {
	chatID := c.Param("chatID")

	var chat Chat
	chatQuery := `SELECT chat_id, users, messages, created_at, updated_at FROM chats WHERE chat_id=$1`
	err := db.QueryRowContext(c, chatQuery, chatID).Scan(&chat.ChatID, pq.Array(&chat.Users), pq.Array(&chat.Messages), &chat.CreatedAt, &chat.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messagesQuery := `SELECT message_id, chat_id, sender, text, media, created_at FROM messages WHERE chat_id=$1 ORDER BY created_at ASC`
	rows, err := db.QueryContext(c, messagesQuery, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.MessageID, &msg.Chat, &msg.Sender, &msg.Text, pq.Array(&msg.Media), &msg.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messages = append(messages, msg)
	}

	c.JSON(http.StatusOK, gin.H{
		"chat":     chat,
		"messages": messages,
	})
}

func GetChatByID(db *sql.DB, c *gin.Context) {
	chatID := c.Param("chatID")

	query := `SELECT chat_id, users, messages, created_at, updated_at FROM chats WHERE chat_id=$1`
	var chat Chat
	err := db.QueryRowContext(c, query, chatID).Scan(&chat.ChatID, pq.Array(&chat.Users), pq.Array(&chat.Messages), &chat.CreatedAt, &chat.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func UpdateChatByID(db *sql.DB, c *gin.Context) {
	chatID := c.Param("chatID")
	var chat Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE chats SET users=$1, messages=$2, updated_at=NOW() WHERE chat_id=$3`
	result, err := db.ExecContext(c, query, pq.Array(chat.Users), pq.Array(chat.Messages), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat updated successfully"})
}

func DeleteChatByID(db *sql.DB, c *gin.Context) {
	chatID := c.Param("chatID")

	query := `DELETE FROM chats WHERE chat_id=$1`
	result, err := db.ExecContext(c, query, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

func DeleteMessageByID(db *sql.DB, c *gin.Context) {
	messageID := c.Param("messageID")

	query := `DELETE FROM messages WHERE message_id=$1`
	result, err := db.ExecContext(c, query, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
