package handler

import (
	"math/rand"
	"net/http"
	"time"

	"whiteboard/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	Hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		Hub: hub,
	}
}

func (h *WebSocketHandler) WebSocketHandler(c *gin.Context) {
	whiteboardID := c.Query("whiteboard_id")
	if whiteboardID == "" {
		whiteboardID = uuid.New().String()
	}

	username := c.Query("username")
	if username == "" {
		username = "user_" + uuid.New().String()[:8]
	}

	color := c.Query("color")
	if color == "" {
		color = randomColor()
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(conn, h.Hub, username, color, whiteboardID)

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func randomColor() string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7",
		"#DDA0DD", "#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E9",
		"#F8B500", "#FF6F61", "#6B5B95", "#88B04B", "#F7CAC9",
	}
	return colors[rand.Intn(len(colors))]
}
