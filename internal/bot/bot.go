package bot

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

type ChatBot struct{}

func NewChatBot() *ChatBot {
	return &ChatBot{}
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
