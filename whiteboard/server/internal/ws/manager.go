package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	MessageTypeJoin   = "join"
	MessageTypeLeave  = "leave"
	MessageTypeCursor = "cursor"
	MessageTypeDraw   = "draw"
	MessageTypeClear  = "clear"
	MessageTypeSync   = "sync"
	MessageTypeUsers  = "users"
)

type Client struct {
	ID           string
	Username     string
	Color        string
	WhiteboardID string
	Conn         *websocket.Conn
	Send         chan []byte
	Hub          *Hub
}

type Message struct {
	Type         string          `json:"type"`
	From         string          `json:"from"`
	WhiteboardID string          `json:"whiteboard_id"`
	Timestamp    int64           `json:"timestamp"`
	Payload      json.RawMessage `json:"payload"`
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Rooms      map[string]map[string]*Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 256),
		Rooms:      make(map[string]map[string]*Client),
	}
}

func NewClient(conn *websocket.Conn, hub *Hub, username, color, whiteboardID string) *Client {
	return &Client{
		ID:           uuid.New().String(),
		Username:     username,
		Color:        color,
		WhiteboardID: whiteboardID,
		Conn:         conn,
		Send:         make(chan []byte, 256),
		Hub:          hub,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Clients[client.ID] = client

	if _, exists := h.Rooms[client.WhiteboardID]; !exists {
		h.Rooms[client.WhiteboardID] = make(map[string]*Client)
	}
	h.Rooms[client.WhiteboardID][client.ID] = client

	h.notifyUserJoin(client)
	h.sendUserList(client.WhiteboardID)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Clients[client.ID]; !exists {
		return
	}

	delete(h.Clients, client.ID)
	close(client.Send)

	if room, exists := h.Rooms[client.WhiteboardID]; exists {
		delete(room, client.ID)
		if len(room) == 0 {
			delete(h.Rooms, client.WhiteboardID)
		}
	}

	h.notifyUserLeave(client)
	h.sendUserList(client.WhiteboardID)
}

func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, exists := h.Rooms[message.WhiteboardID]
	if !exists {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, client := range room {
		if client.ID == message.From && message.Type != MessageTypeSync {
			continue
		}
		select {
		case client.Send <- data:
		default:
			h.mu.RUnlock()
			h.mu.Lock()
			if c, ok := h.Clients[client.ID]; ok {
				close(c.Send)
				delete(h.Clients, client.ID)
				if room, exists := h.Rooms[c.WhiteboardID]; exists {
					delete(room, client.ID)
					if len(room) == 0 {
						delete(h.Rooms, c.WhiteboardID)
					}
				}
			}
			h.mu.Unlock()
			h.mu.RLock()
		}
	}
}

func (h *Hub) notifyUserJoin(client *Client) {
	payload := map[string]interface{}{
		"user_id":   client.ID,
		"username":  client.Username,
		"color":     client.Color,
	}
	payloadBytes, _ := json.Marshal(payload)

	msg := &Message{
		Type:         MessageTypeJoin,
		From:         client.ID,
		WhiteboardID: client.WhiteboardID,
		Timestamp:    time.Now().UnixMilli(),
		Payload:      payloadBytes,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	room, exists := h.Rooms[client.WhiteboardID]
	if !exists {
		return
	}

	for _, c := range room {
		if c.ID == client.ID {
			continue
		}
		select {
		case c.Send <- data:
		default:
		}
	}
}

func (h *Hub) notifyUserLeave(client *Client) {
	payload := map[string]interface{}{
		"user_id":   client.ID,
		"username":  client.Username,
		"color":     client.Color,
	}
	payloadBytes, _ := json.Marshal(payload)

	msg := &Message{
		Type:         MessageTypeLeave,
		From:         client.ID,
		WhiteboardID: client.WhiteboardID,
		Timestamp:    time.Now().UnixMilli(),
		Payload:      payloadBytes,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	room, exists := h.Rooms[client.WhiteboardID]
	if !exists {
		return
	}

	for _, c := range room {
		select {
		case c.Send <- data:
		default:
		}
	}
}

func (h *Hub) sendUserList(whiteboardID string) {
	room, exists := h.Rooms[whiteboardID]
	if !exists {
		return
	}

	users := make([]map[string]interface{}, 0, len(room))
	for _, c := range room {
		users = append(users, map[string]interface{}{
			"user_id":  c.ID,
			"username": c.Username,
			"color":    c.Color,
		})
	}

	payload := map[string]interface{}{
		"users": users,
	}
	payloadBytes, _ := json.Marshal(payload)

	msg := &Message{
		Type:         MessageTypeUsers,
		From:         "system",
		WhiteboardID: whiteboardID,
		Timestamp:    time.Now().UnixMilli(),
		Payload:      payloadBytes,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for _, c := range room {
		select {
		case c.Send <- data:
		default:
		}
	}
}
