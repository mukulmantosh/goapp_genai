package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
)

const (
	testWebSocketURL = "ws://localhost:8080/ws/model"
)

func TestWebSocketConnection(t *testing.T) {
	// Create WebSocket URL with query parameters
	u, err := url.Parse(testWebSocketURL)
	if err != nil {
		t.Fatalf("Failed to parse WebSocket base URL: %v", err)
	}
	q := u.Query()
	q.Set("model", "llama3")
	u.RawQuery = q.Encode()

	// Dial the WebSocket server
	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket server: %v", err)
	}
	defer client.Close()

	message := "Tell me something about yourself in 20 words."
	if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	_, receivedMessage, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}
	log.Printf("Received message: %s", string(receivedMessage))
	if len(string(receivedMessage)) <= 0 {
		t.Error("Received empty message")
	}
}
