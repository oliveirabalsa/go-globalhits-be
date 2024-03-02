package main

import (
	"encoding/json"
	"log"

	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
	"github.com/oliveirabalsa/go-globalhitss-be/config"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
)

func main() {
	ch, conn, db := config.InitServices()
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

	clientRepo := repository.NewClientRepository(db)

	clientUsecase := usecase.NewClientUseCase(*clientRepo)

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

				_, err = clientUsecase.CreateClient(&client)
				if err != nil {
					log.Printf("Failed to create client: %v", err)
					continue
				}
			}
		}
	}()

	log.Printf("Worker running. To exit press CTRL+C")
	<-forever
}
