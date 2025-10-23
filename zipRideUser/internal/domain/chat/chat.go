package chat

import (
	"log"
	"net/http"
	"strconv"
	"zipride/database"
	"zipride/internal/models"
	connection "zipride/internal/websocket"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)
//upgrader for the websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
//function for connect user
func ConnectUser(c *gin.Context) {
	// Get booking ID from path parameter
	bookingIDStr := c.Param("booking_id")
	if bookingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking_id is required"})
		return
	}

	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_id"})
		return
	}

	// Get user ID from middleware
	currentUserID := middleware.GetUserID(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated user"})
		return
	}

	// Fetch booking from database
	var booking models.Booking
	if err := database.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "booking not found"})
		return
	}

	// Verify the user owns this booking
	if booking.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized for this booking"})
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
		Role:      "user",
	}

	connection.AddClient(client)
	go connection.HandleClientMessages(client)

	// Send welcome message
	conn.WriteJSON(models.ChatMessage{
		BookingID: uint(bookingID),
		Sender:    "system",
		Message:   "Connected as user",
	})
}
