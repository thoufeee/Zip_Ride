package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"zipride/internal/models"

	"github.com/segmentio/kafka-go"
)

func Producer(booking models.BookingMessage) error {
	// Convert the booking struct to JSON
	data, err := json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}
	// creates a producer connection to Kafka.
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"}, //Kafka broker addresses.
		Topic:    "booking-events",           //topic to send messages to
		Balancer: &kafka.LeastBytes{},        //decides how to distribute messages among partitions
	})
	defer writer.Close() //connection closes automatically when the function ends.

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", booking.BookingID)), //A unique identifier
		Value: data,                                         //actual data
	}

	time.Sleep(500 * time.Millisecond) // optional sync delay
	// Send Message to Kafka
	if err := writer.WriteMessages(context.Background(), msg); err != nil {
		log.Println("Kafka write error:", err)
		return err
	}
	//Log Success
	fmt.Printf("✅ Kafka Message Sent — BookingID: %d | Vehicle: %s | Fare: %.2f\n", booking.BookingID, booking.Vehicle, booking.Fare)
	return nil
}