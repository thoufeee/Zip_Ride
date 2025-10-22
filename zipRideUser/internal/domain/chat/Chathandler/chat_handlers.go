package chathandler

import (
	"log"
	"net/http"
	"strconv"
	"zipride/database"
	"zipride/internal/domain/chat/services"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader config: allows all origins
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ChatWebSocket handles WebSocket connections for users only
func ChatWebSocket(c *gin.Context) {
	// Get booking_id from query parameter
	bookingIDStr := c.Query("booking_id")
	if bookingIDStr == "" {
		// Return error if booking_id is missing
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking_id is required"})
		return
	}

	// Convert booking_id from string to uint64
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 64)
	if err != nil {
		// Return error if booking_id is invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
		return
	}

	// Get user ID from JWT middleware (user must be authenticated)
	userID := c.GetUint("user_id")

	// Verify that the booking exists and belongs to the logged-in user
	var booking models.Booking
	if err := database.DB.First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		// Return forbidden if booking not found or user is not the owner
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized or booking not found"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Create a client object representing this user connection
	client := &services.Client{
		Conn:      conn,
		BookingID: uint(bookingID),
		Role:      "user",
	}

	// Register the client in the global clients map
	services.AddClient(client)

	// Start a goroutine to handle incoming messages from this client
	go services.HandleClientMessages(client)

	// Send a welcome message to the connected user
	welcome := models.ChatMessage{
		BookingID: uint(bookingID),
		Sender:    "system",
		Message:   "Secure user chat connected",
	}
	conn.WriteJSON(welcome)
}
