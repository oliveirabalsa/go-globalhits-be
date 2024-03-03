package usecase

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/dto"
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
		Action: os.Getenv("CREATE_CLIENT_ACTION"),
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
		Action: os.Getenv("UPDATE_CLIENT_ACTION"),
		Data:   clientData,
	}
	err = queue.PublishMessage(uc.QueueChannel, uc.QueueName, message)
	if err != nil {
		return "", err
	}
	return "Your data has been received and is being processed.", nil
}

func (uc *ClientUsecase) GetClients(page int, pageSize int) (dto.PaginationResponse, error) {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	clients, totalPages, err := uc.ClientRepo.GetAll(page, pageSize)
	if err != nil {
		return dto.PaginationResponse{}, err
	}

	for _, client := range clients {
		client.DecryptSensitiveData()
	}

	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	response := dto.PaginationResponse{
		Data:       clients,
		Page:       page,
		NextPage:   nextPage,
		TotalPages: totalPages,
	}

	return response, nil
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
		Action: os.Getenv("DELETE_CLIENT_ACTION"),
		Data:   []byte(clientId.String()),
	}
	if err := queue.PublishMessage(uc.QueueChannel, uc.QueueName, message); err != nil {
		return "", err
	}

	return "Your data has been received and is being processed.", nil
}
