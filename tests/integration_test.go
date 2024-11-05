package tests

import (
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketStream(t *testing.T) {
	wavData, err := os.ReadFile("../harvard.wav")
	if err != nil {
		t.Fatalf("Failed to read test WAV file: %v", err)
	}

	// Create a WebSocket connection
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/stream"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send WAV data to the server
	err = ws.WriteMessage(websocket.BinaryMessage, wavData)
	if err != nil {
		t.Fatalf("Failed to send WAV data: %v", err)
	}

	_, flacData, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to receive FLAC data: %v", err)
	}

	if len(flacData) == 0 {
		t.Error("FLAC data is empty")
	}
}
