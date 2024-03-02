package main

import (
	"encoding/json"
	"log"

	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
)

func main() {
	conn, ch, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	queueName := "globalhitss"
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var message queue.Message
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			switch message.Action {
			case "create_client":
				var client model.Client
				err := json.Unmarshal(message.Data, &client)
				if err != nil {
					log.Printf("Failed to unmarshal client data: %v", err)
					continue
				}

				log.Printf("Received message to create client: %+v", client)
			}
		}
	}()

	log.Printf("Worker running. To exit press CTRL+C")
	<-forever
}
