// database.go
package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// Database represents the SQLite database connection.
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection and initializes the necessary tables.
func NewDatabase(databasePath string) (*Database, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY,
            username TEXT
        );

        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY,
            sender_id INTEGER,
            room_name TEXT,
            content TEXT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	return d.db.Close()
}

// AddUser adds a new user to the database.
func (d *Database) AddUser(username string) (int64, error) {
	result, err := d.db.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// AddMessage adds a new message to the database.
func (d *Database) AddMessage(senderID int64, roomName, content string) (int64, error) {
	result, err := d.db.Exec("INSERT INTO messages (sender_id, room_name, content) VALUES (?, ?, ?)", senderID, roomName, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetRoomMessages retrieves the last 'limit' messages for a specific room.
func (d *Database) GetRoomMessages(roomName string, limit int) ([]Message, error) {
	rows, err := d.db.Query("SELECT sender_id, content, timestamp FROM messages WHERE room_name = ? ORDER BY timestamp DESC LIMIT ?", roomName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var senderID int
		var content string
		var timestamp string
		err := rows.Scan(&senderID, &content, &timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, Message{SenderID: senderID, Content: content, Timestamp: timestamp})
	}
	return messages, nil
}

// Message represents a chat message retrieved from the database.
type Message struct {
	SenderID  int
	Content   string
	Timestamp string
}
