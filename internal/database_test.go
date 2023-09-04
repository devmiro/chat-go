// database_test.go
package main

import (
	"database/sql"
	"os"
	"testing"
)

const testDBPath = "test.db"

func TestMain(m *testing.M) {
	// Set up the test database before running tests
	if err := createTestDB(); err != nil {
		os.Exit(1)
	}
	code := m.Run()
	// Clean up and remove the test database after running tests
	if err := cleanUpTestDB(); err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}

func createTestDB() error {
	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		return err
	}
	defer db.Close()

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
	return err
}

func cleanUpTestDB() error {
	return os.Remove(testDBPath)
}

func TestDatabase_AddUser(t *testing.T) {
	db, err := NewDatabase(testDBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	username := "testuser"
	_, err = db.AddUser(username)
	if err != nil {
		t.Errorf("AddUser() error = %v, want nil", err)
	}
}

func TestDatabase_AddMessage(t *testing.T) {
	db, err := NewDatabase(testDBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	senderID := int64(1)
	roomName := "testroom"
	content := "Hello, world!"

	_, err = db.AddMessage(senderID, roomName, content)
	if err != nil {
		t.Errorf("AddMessage() error = %v, want nil", err)
	}
}

func TestDatabase_GetRoomMessages(t *testing.T) {
	db, err := NewDatabase(testDBPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	roomName := "testroom"
	limit := 10

	messages, err := db.GetRoomMessages(roomName, limit)
	if err != nil {
		t.Errorf("GetRoomMessages() error = %v, want nil", err)
	}

	// Check the number of retrieved messages
	if len(messages) != 0 {
		t.Errorf("GetRoomMessages() got %d messages, want 0", len(messages))
	}
}
