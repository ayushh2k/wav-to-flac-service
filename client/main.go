// client/main.go

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	wavFile  = flag.String("wav", "test.wav", "path to the WAV file")
	flacFile = flag.String("flac", "output.flac", "path to save the FLAC file")
)

func main() {
	flag.Parse()

	// Open the WAV file
	wavData, err := os.ReadFile(*wavFile)
	if err != nil {
		log.Fatalf("Failed to read WAV file: %v", err)
	}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+*addr+"/stream", nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Send the WAV data
	err = conn.WriteMessage(websocket.BinaryMessage, wavData)
	if err != nil {
		log.Fatalf("Failed to send WAV data: %v", err)
	}

	// Set a timeout for receiving the FLAC data
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// Receive the FLAC data
	_, flacData, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("Failed to receive FLAC data: %v", err)
	}

	// Save the FLAC data to a file
	err = os.WriteFile(*flacFile, flacData, 0644)
	if err != nil {
		log.Fatalf("Failed to save FLAC file: %v", err)
	}

	log.Println("FLAC file saved successfully:", *flacFile)
}
