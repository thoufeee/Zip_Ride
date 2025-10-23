package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"zipride/internal/models"

	"github.com/segmentio/kafka-go"
)

func Consumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "booking-events",
		GroupID: "booking-group",
	})
	defer reader.Close()	

	fmt.Println("👂 Kafka Consumer listening on topic: booking-events")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var booking models.BookingMessage
		if err := json.Unmarshal(msg.Value, &booking); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		fmt.Printf("🚖 Booking Received → ID: %d | UserID: %d | Vehicle: %s | Fare: %.2f | Status: %s\n",
			booking.BookingID, booking.UserID, booking.Vehicle, booking.Fare, booking.Status)
	}
}
