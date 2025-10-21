package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"zipride/internal/models"

	"github.com/IBM/sarama"
)

var bookingProducer sarama.SyncProducer
var bookingTopic string

// InitBookingProducer initializes Kafka producer
func InitBookingProducer() {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	bookingTopic = os.Getenv("KAFKA_BOOKING_TOPIC")
	if bookingTopic == "" {
		bookingTopic = "booking-events"
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.Timeout = 5 * time.Second

	var err error
	bookingProducer, err = sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		log.Fatalf("Kafka producer init failed: %v", err)
	}

	log.Println("Kafka booking producer initialized")
}

// PublishBookingEvent sends booking event to Kafka
func PublishBookingEvent(booking models.Booking) error {
	if bookingProducer == nil {
		return fmt.Errorf("Kafka producer not initialized")
	}

	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: bookingTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", booking.UserID)),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := bookingProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Booking event sent: partition=%d offset=%d\n", partition, offset)
	return nil
}
