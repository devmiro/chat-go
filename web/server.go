package web

import (
	"chat-go/internal/bot"
	"chat-go/internal/message"
	"chat-go/internal/user"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ChatServer struct {
	chatRooms map[string]*ChatRoom
	upgrader  websocket.Upgrader
	// Add other dependencies like user manager, message service, and chat bot
	userManager    *user.UserManager
	messageService *message.MessageService
	bot            *bot.ChatBot
}

type ChatRoom struct {
	Name     string
	Messages []*message.Message
	Clients  map[*websocket.Conn]bool
}

func NewChatServer(userManager *user.UserManager, messageService *message.MessageService, chatBot *bot.ChatBot) *ChatServer {
	return &ChatServer{
		chatRooms:      map[string]*ChatRoom{},
		upgrader:       websocket.Upgrader{},
		userManager:    userManager,
		messageService: messageService,
		bot:            chatBot,
	}
}

// func NewChatServer() *ChatServer {
// 	return &ChatServer{
// 		chatRooms: make(map[string]*ChatRoom),
// 		upgrader: websocket.Upgrader{
// 			ReadBufferSize:  1024,
// 			WriteBufferSize: 1024,
// 		},
// 	}
// }

func (s *ChatServer) handleChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomName := vars["roomName"]

	if _, ok := s.chatRooms[roomName]; !ok {
		// Create a new chat room if it doesn't exist
		s.chatRooms[roomName] = &ChatRoom{
			Name:     roomName,
			Messages: make([]*message.Message, 0),
			Clients:  make(map[*websocket.Conn]bool),
		}
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	room := s.chatRooms[roomName]
	room.Clients[conn] = true

	// Handle incoming messages from clients
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		messageContent := string(p)

		// Handle commands, e.g., /stock=stock_code
		if messageContent[0] == '/' {
			response, err := s.bot.RespondToCommand(messageContent, s.messageService)
			if err != nil {
				log.Println(err)
				return
			}
			// Broadcast the bot's response to all clients
			for client := range room.Clients {
				err := client.WriteMessage(messageType, []byte(response))
				if err != nil {
					log.Println(err)
					return
				}
			}
		} else {
			// Handle regular chat messages
			sender := s.userManager.RegisterUser("Anonymous") // You can implement user authentication here
			msg := s.messageService.SendMessage(sender.ID, room.Name, messageContent)

			// Broadcast the message to all clients
			for client := range room.Clients {
				// Send the message with sender's username
				msgWithUsername := fmt.Sprintf("%s: %s", sender.Username, msg.Content)
				err := client.WriteMessage(messageType, []byte(msgWithUsername))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func (s *ChatServer) getChatMessages(roomName string, limit int) []*message.Message {
	// Get recent chat messages for a room
	messages := s.messageService.GetRoomMessages(roomName, limit)
	return messages
}

func (s *ChatServer) renderChatPage(w http.ResponseWriter, r *http.Request, roomName string) {
	tpl := template.Must(template.ParseFiles("web/templates/index.html"))

	messages := s.getChatMessages(roomName, 50) // Retrieve the last 50 messages
	tplData := struct {
		RoomName string
		Messages []*message.Message
	}{
		RoomName: roomName,
		Messages: messages,
	}

	if err := tpl.Execute(w, tplData); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *ChatServer) serveStaticFiles() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
}

func (s *ChatServer) serveHTMLTemplates() {
	tpl := template.Must(template.ParseFiles("D:/workspace/chat-go/web/templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomName := vars["roomName"]
		tpl.Execute(w, roomName)
	})
}

func (s *ChatServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/chat/{roomName}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomName := vars["roomName"]
		s.renderChatPage(w, r, roomName)
	})
	r.HandleFunc("/chatsocket/{roomName}", s.handleChat)

	s.serveStaticFiles()
	s.serveHTMLTemplates()

	//http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
