package bot

import (
	"chat-go/internal/message"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

type ChatBot struct{}

func NewChatBot() *ChatBot {
	return &ChatBot{}
}

func (b *ChatBot) RespondToCommand(command string, roomName string, messageService *message.MessageService) (string, error) {
	// Check if the message is a command
	if strings.HasPrefix(command, "/stock=") {
		// Extract the stock code from the command
		stockCode := strings.TrimPrefix(command, "/stock=")

		// Call the GetStockQuote method to fetch the stock quote
		quote, err := b.GetStockQuote(stockCode)
		if err != nil {
			return "", err
		}

		// Create a response message
		responseMessage := "StockBot" + ": " + quote

		// Post the response message to the chatroom using messageService
		err = messageService.SendMessage(99, roomName, responseMessage)
		if err != nil {
			return "", err
		}

		return responseMessage, nil
	}

	// If it's not a recognized command, return an empty response
	return "", nil
}

func (b *ChatBot) GetStockQuote(stockCode string) (string, error) {
	// Perform API request to retrieve stock data
	// Parse CSV response and extract the stock quote
	// Example:
	response, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	reader := csv.NewReader(response.Body)
	record, err := reader.Read()
	if err != nil {
		return "", err
	}

	// Extract relevant data from the CSV
	stockQuote := fmt.Sprintf("%s quote is $%s per share", stockCode, record[6])
	return stockQuote, nil
}
