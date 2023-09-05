package message

import (
	"fmt"
	"time"
)

type Message struct {
	ID        int
	SenderID  int
	RoomName  string
	Content   string
	Timestamp time.Time
}

// Error implements error.
func (*Message) Error() string {
	fmt.Println("Message error: Message is not valid")
	return "Message error: Message is not valid"
}

type MessageService struct {
	Messages []*Message
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (ms *MessageService) SendMessage(senderID int, roomName, content string) *Message {
	message := &Message{
		ID:        len(ms.Messages) + 1,
		SenderID:  senderID,
		RoomName:  roomName,
		Content:   content,
		Timestamp: time.Now(),
	}

	ms.Messages = append(ms.Messages, message)
	return message
}

func (ms *MessageService) GetRoomMessages(roomName string, limit int) []*Message {
	var roomMessages []*Message
	for i := len(ms.Messages) - 1; i >= 0; i-- {
		if ms.Messages[i].RoomName == roomName {
			roomMessages = append(roomMessages, ms.Messages[i])
			if len(roomMessages) >= limit {
				break
			}
		}
	}
	return roomMessages
}
