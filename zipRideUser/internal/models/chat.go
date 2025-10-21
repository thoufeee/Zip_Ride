package models

//struct for chatting messages
type ChatMessage struct {
	BookingID uint   `json:"booking_id"`
	Sender    string `json:"sender"` // "user" or "driver"
	Message   string `json:"message"`
}