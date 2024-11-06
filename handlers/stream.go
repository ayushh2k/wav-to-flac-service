// handlers/stream.go
package handlers

import (
	"io"
	"net/http"
	"wav-to-flac-service/services"
	"wav-to-flac-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Size of audio chunks to process at once
	chunkSize = 4096
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StreamHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.LogError("upgrade", err)
		return
	}
	defer ws.Close()

	// Create conversion pipeline
	converter, err := services.NewStreamingConverter()
	if err != nil {
		utils.LogError("converter creation", err)
		return
	}
	defer converter.Close()

	// Create channels for communication
	errChan := make(chan error, 1)
	dataChan := make(chan []byte, 100)

	// Start goroutine to read converted FLAC data
	go func() {
		buffer := make([]byte, chunkSize)
		for {
			n, err := converter.ReadConverted(buffer)
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				close(dataChan)
				return
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])
				dataChan <- data
			}
		}
	}()

	// Start goroutine to send converted data back to client
	go func() {
		for data := range dataChan {
			if err := ws.WriteMessage(websocket.BinaryMessage, data); err != nil {
				errChan <- err
				return
			}
		}
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.LogError("read", err)
			}
			break
		}

		if err := converter.WriteInput(message); err != nil {
			utils.LogError("write to converter", err)
			break
		}
	}
}
