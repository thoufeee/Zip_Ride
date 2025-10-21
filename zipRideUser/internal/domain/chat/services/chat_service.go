package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"zipride/internal/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	BookingID uint
	Role      string
}

var (
	clients = make(map[uint][]*Client) // bookingID -> clients
	mutex   = &sync.Mutex{}
)

// Add a new client
func AddClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
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
}

// Broadcast message to all clients for the booking
func BroadcastMessage(msg models.ChatMessage) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range clients[msg.BookingID] {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("Write error:", err)
		}
	}
}

// Delete all chats after booking is completed
func DeleteBookingChats(bookingID uint) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range clients[bookingID] {
		client.Conn.Close()
	}
	delete(clients, bookingID)
}

// Handle incoming messages for a client
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

		BroadcastMessage(msg)
	}
}

// Utility to convert string bookingID to uint
func StringToUint(idStr string) (uint, error) {
	var id uint
	_, err := fmt.Sscan(idStr, &id)
	return id, err
}
