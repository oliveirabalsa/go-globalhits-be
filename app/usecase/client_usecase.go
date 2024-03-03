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
	client.EncryptSensitiveData()
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
	client.EncryptSensitiveData()
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

func (uc *ClientUsecase) GetClients(page int, pageSize int) ([]*model.Client, int, error) {
	clients, totalPages, err := uc.ClientRepo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	for _, client := range clients {
		client.DecryptSensitiveData()
	}

	return clients, totalPages, nil
}

func (uc *ClientUsecase) GetClientByID(clientId uuid.UUID) (*model.Client, error) {
	client, err := uc.ClientRepo.GetByID(clientId)
	if err != nil {
		return nil, err
	}

	client.DecryptSensitiveData()
	return client, nil
}

func (uc *ClientUsecase) DeleteClient(clientId uuid.UUID) (string, error) {
	message := &queue.Message{
		Action: "delete_client",
		Data:   []byte(clientId.String()),
	}
	if err := queue.PublishMessage(uc.QueueChannel, uc.QueueName, message); err != nil {
		return "", err
	}

	return "Your data has been received and is being processed.", nil
}
