// handlers/stream.go

package handlers

import (
	"log"
	"net/http"

	"wav-to-flac-service/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StreamHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		// Read WAV data from WebSocket
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		// Convert WAV to FLAC using service
		flacData, err := services.WavToFlac(message)
		if err != nil {
			log.Println("wavToFlac:", err)
			continue
		}

		// Send FLAC data back to client
		err = ws.WriteMessage(websocket.BinaryMessage, flacData)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
