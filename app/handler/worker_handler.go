package handler

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
)

func CreateClientHandler(clientRepo *repository.ClientRepository, message queue.Message) {
	var client model.Client
	err := json.Unmarshal(message.Data, &client)
	if err != nil {
		log.Printf("Failed to unmarshal client data: %v", err)
		return
	}

	_, err = clientRepo.Create(&client)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return
	}
}

func UpdateClientHandler(clientRepo *repository.ClientRepository, message queue.Message) {
	var client model.Client
	err := json.Unmarshal(message.Data, &client)
	if err != nil {
		log.Printf("Failed to unmarshal client data: %v", err)
		return
	}

	_, err = clientRepo.Update(&client)
	if err != nil {
		log.Printf("Failed to update client: %v", err)
		return
	}
}

func DeleteClientHandler(clientRepo *repository.ClientRepository, message queue.Message) {
	clientId, err := uuid.Parse(string(message.Data))
	if err != nil {
		log.Printf("Failed to parse client ID: %v", err)
		return
	}

	err = clientRepo.DeleteClient(clientId)
	if err != nil {
		log.Printf("Failed to delete client: %v", err)
		return
	}
}
