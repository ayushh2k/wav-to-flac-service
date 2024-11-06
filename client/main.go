// client/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Default chunk size for reading WAV file (32KB)
	defaultChunkSize = 32 * 1024
)

type Client struct {
	conn       *websocket.Conn
	url        string
	chunkSize  int
	outputDir  string
	logger     *log.Logger
	receivedWg sync.WaitGroup
}

func NewClient(serverURL, outputDir string, chunkSize int) *Client {
	return &Client{
		url:       serverURL,
		chunkSize: chunkSize,
		outputDir: outputDir,
		logger:    log.New(os.Stdout, "[Client] ", log.LstdFlags),
	}
}

func (c *Client) Connect() error {
	u, err := url.Parse(c.url)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	c.logger.Printf("Connecting to %s", c.url)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}

	c.conn = conn
	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) ProcessFile(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	outputPath := filepath.Join(c.outputDir, fmt.Sprintf("%s.flac",
		filepath.Base(inputFile[:len(inputFile)-len(filepath.Ext(inputFile))])))

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Start goroutine to receive converted FLAC data
	c.receivedWg.Add(1)
	go func() {
		defer c.receivedWg.Done()
		for {
			messageType, data, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.logger.Println("Connection closed normally")
					return
				}
				c.logger.Printf("Read error: %v", err)
				return
			}

			if messageType == websocket.BinaryMessage {
				_, err := outFile.Write(data)
				if err != nil {
					c.logger.Printf("Failed to write FLAC data: %v", err)
					return
				}
			} else if messageType == websocket.TextMessage {
				c.logger.Printf("Server message: %s", string(data))
			}
		}
	}()

	// Read and send WAV file in chunks
	buffer := make([]byte, c.chunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read file: %v", err)
		}

		err = c.conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			return fmt.Errorf("failed to send data: %v", err)
		}

		// Small delay to prevent overwhelming the server
		time.Sleep(10 * time.Millisecond)
	}

	err = c.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		c.logger.Printf("Error sending close message: %v", err)
	}

	// Wait for all data to be received
	c.receivedWg.Wait()
	c.logger.Printf("Conversion completed. Output saved to: %s", outputPath)
	return nil
}

func main() {
	serverURL := flag.String("server", "ws://localhost:8080/stream",
		"WebSocket server URL")
	inputFile := flag.String("input", "", "Input WAV file path")
	outputDir := flag.String("output", ".", "Output directory for FLAC files")
	chunkSize := flag.Int("chunk", defaultChunkSize, "Chunk size in bytes")

	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Input file is required")
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	client := NewClient(*serverURL, *outputDir, *chunkSize)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	if err := client.ProcessFile(*inputFile); err != nil {
		log.Fatalf("Failed to process file: %v", err)
	}
}
