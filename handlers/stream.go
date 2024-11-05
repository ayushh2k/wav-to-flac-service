// handlers/stream.go

package handlers

import (
	"net/http"

	"wav-to-flac-service/services"
	"wav-to-flac-service/utils"

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
		utils.LogError("upgrade", err)
		return
	}
	defer ws.Close()

	for {
		// Read WAV data from WebSocket
		_, message, err := ws.ReadMessage()
		if err != nil {
			utils.LogError("read", err)
			break
		}

		// Convert WAV to FLAC using service
		flacData, err := services.WavToFlac(message)
		if err != nil {
			utils.LogError("wavToFlac", err)
			continue
		}

		// Send FLAC data back to client
		err = ws.WriteMessage(websocket.BinaryMessage, flacData)
		if err != nil {
			utils.LogError("write", err)
			break
		}
	}
}
