package main

import (
	"fmt"

	"chat-go/internal/bot"
	"chat-go/internal/message"
	"chat-go/internal/user"
	"chat-go/web"
)

func main() {
	// Create instances of the required components
	userManager := user.NewUserManager()
	messageService := message.NewMessageService()
	chatBot := bot.NewChatBot()

	// Create and start the chat server
	chatServer := web.NewChatServer(userManager, messageService, chatBot) // Adjust the dependencies
	go chatServer.Start()

	fmt.Println("Chat server is running on port 8080...")

	// Handle server shutdown gracefully (you can add more graceful shutdown logic)
	select {}
}
