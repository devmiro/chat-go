package chat

import (
	"sync"
	"time"
)

type Message struct {
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatRoom struct {
	ID       int
	Name     string
	Messages []Message
	mu       sync.Mutex
}

func NewChatRoom(name string) *ChatRoom {
	return &ChatRoom{
		Name: name,
	}
}

func (c *ChatRoom) AddMessage(userID int, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	message := Message{
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
	}

	c.Messages = append(c.Messages, message)
	// Ensure that only the last 50 messages are stored
	if len(c.Messages) > 50 {
		c.Messages = c.Messages[len(c.Messages)-50:]
	}
}

func (c *ChatRoom) GetMessages() []Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Messages
}
