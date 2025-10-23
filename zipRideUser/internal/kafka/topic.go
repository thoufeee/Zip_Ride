package kafka

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func CreateTopic() {
	//create connection
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal("Failed to connect to Kafka:", err)
	}
	defer conn.Close() //ensure connection close automaticlly
	//create controller
	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Failed to get controller:", err)
	}
	//create connection to the controller
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatal("Failed to connect to controller:", err)
	}
	defer controllerConn.Close() //ensure connection close automaticlly

	topic := "booking-events" //the topic name

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	//if something fails
	if err != nil {
		log.Println("⚠️ Topic might already exist:", err)
	} else {
		//sucess responce
		fmt.Println("✅ Kafka topic created:", topic)
	}
}
