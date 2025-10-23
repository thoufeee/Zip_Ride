package chat

import (
	"log"
	"net/http"
	"strconv"
	"zipride/database"
	"zipride/internal/models"
	connection "zipride/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ConnectUserOrDriver handles both user and driver
func ConnectUser(c *gin.Context) {
	bookingIDStr := c.Query("booking_id")
	if bookingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking_id is required"})
		return
	}

	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
		return
	}

	role := c.Query("role") // "user" or "driver"
	if role != "user" && role != "driver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be user or driver"})
		return
	}

	// Optional: validate user/driver belongs to the booking
	var booking models.Booking
	if err := database.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "booking not found"})
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	client := &connection.Client{
		Conn:      conn,
		BookingID: uint(bookingID),
		Role:      role,
	}

	connection.AddClient(client)
	go connection.HandleClientMessages(client)

	// Welcome message
	conn.WriteJSON(models.ChatMessage{
		BookingID: uint(bookingID),
		Sender:    "system",
		Message:   "Connected as " + role,
	})
}
