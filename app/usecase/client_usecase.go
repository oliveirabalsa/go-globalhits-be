package usecase

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
	"github.com/streadway/amqp"
)

type ClientUsecase struct {
	ClientRepo   repository.ClientRepository
	QueueChannel *amqp.Channel
	QueueName    string
}

func NewClientUseCase(clientRepo repository.ClientRepository, queueChannel *amqp.Channel, queueName string) *ClientUsecase {
	return &ClientUsecase{
		ClientRepo:   clientRepo,
		QueueChannel: queueChannel,
		QueueName:    queueName,
	}
}

func (uc *ClientUsecase) CreateClient(client *model.Client) (string, error) {
	clientData, err := json.Marshal(client)
	if err != nil {
		return "", err

	}
	message := &queue.Message{
		Action: "create_client",
		Data:   clientData,
	}
	err = queue.PublishMessage(uc.QueueChannel, uc.QueueName, message)
	if err != nil {
		return "", err
	}
	return "Your data has been received and is being processed.", nil
}

func (uc *ClientUsecase) UpdateClient(clientId uuid.UUID, client *model.Client) (string, error) {
	client.ID = clientId
	clientData, err := json.Marshal(client)
	if err != nil {
		return "", err

	}

	message := &queue.Message{
		Action: "update_client",
		Data:   clientData,
	}
	err = queue.PublishMessage(uc.QueueChannel, uc.QueueName, message)
	if err != nil {
		return "", err
	}
	return "Your data has been received and is being processed.", nil
}

func (uc *ClientUsecase) GetClients() ([]*model.Client, error) {
	return uc.ClientRepo.GetAll()
}

func (uc *ClientUsecase) DeleteClient(clientId uuid.UUID) (string, error) {
	if err := uc.ClientRepo.ClientExists(clientId); err != nil {
		return "", err
	}

	message := &queue.Message{
		Action: "delete_client",
		Data:   []byte(clientId.String()),
	}
	if err := queue.PublishMessage(uc.QueueChannel, uc.QueueName, message); err != nil {
		return "", err
	}

	return "Your data has been received and is being processed.", nil
}
