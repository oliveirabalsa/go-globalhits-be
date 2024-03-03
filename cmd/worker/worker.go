package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/oliveirabalsa/go-globalhitss-be/app/handler"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/config"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	ch, conn, db := config.InitServices()
	defer conn.Close()
	defer ch.Close()

	msgs, err := config.SetupConsumer(ch)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	clientRepo := repository.NewClientRepository(db)
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
			case os.Getenv("CREATE_CLIENT_ACTION"):
				handler.CreateClientHandler(clientRepo, message)
			case os.Getenv("UPDATE_CLIENT_ACTION"):
				handler.UpdateClientHandler(clientRepo, message)
			case os.Getenv("DELETE_CLIENT_ACTION"):
				handler.DeleteClientHandler(clientRepo, message)
			}
		}
	}()

	log.Printf("Starting worker...")
	<-forever
}
