package connection

import (
	"encoding/json"
	"log"
	"sync"
	"zipride/internal/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	BookingID uint
	Role      string // "user" or "driver"
}

var (
	clients = make(map[uint][]*Client) // bookingID -> clients (user + driver)
	mutex   = &sync.Mutex{}
)

// Add a new client
func AddClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	// Only allow max 2 clients per booking (user + driver)
	if len(clients[client.BookingID]) >= 2 {
		client.Conn.WriteJSON(models.ChatMessage{
			BookingID: client.BookingID,
			Sender:    "system",
			Message:   "Chat full: only one user and one driver allowed",
		})
		client.Conn.Close()
		return
	}

	clients[client.BookingID] = append(clients[client.BookingID], client)
}

// Remove client
func RemoveClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	list := clients[client.BookingID]
	for i, c := range list {
		if c == client {
			clients[client.BookingID] = append(list[:i], list[i+1:]...)
			break
		}
	}
	if len(clients[client.BookingID]) == 0 {
		delete(clients, client.BookingID)
	}
}

// Broadcast message to the other client in the booking
func BroadcastMessage(msg models.ChatMessage, sender *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range clients[msg.BookingID] {
		if client != sender { // send only to the other client
			if err := client.Conn.WriteJSON(msg); err != nil {
				log.Println("Write error:", err)
			}
		}
	}
}

// Handle incoming messages
func HandleClientMessages(client *Client) {
	defer func() {
		client.Conn.Close()
		RemoveClient(client)
	}()

	for {
		_, msgBytes, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg models.ChatMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		// Broadcast to the other client
		BroadcastMessage(msg, client)
	}
}
