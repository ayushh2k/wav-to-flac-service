// tests/streaming_test.go
package tests

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestStreamingConversion(t *testing.T) {
	// Open test WAV file
	file, err := os.Open("../harvard.wav")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// Connect to WebSocket
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/stream"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Create channels for test coordination
	done := make(chan struct{})
	flacData := new(bytes.Buffer)

	// Start goroutine to receive converted data
	go func() {
		defer close(done)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				if err != io.EOF {
					t.Errorf("Error reading WebSocket: %v", err)
				}
				return
			}
			flacData.Write(message)
		}
	}()

	// Stream WAV data in chunks
	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading WAV file: %v", err)
		}

		err = ws.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			t.Fatalf("Error writing to WebSocket: %v", err)
		}

		// Small delay to simulate real-time streaming
		time.Sleep(time.Millisecond * 10)
	}

	// Wait for processing to complete
	select {
	case <-done:
		if flacData.Len() == 0 {
			t.Error("No FLAC data received")
		}
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}
}
