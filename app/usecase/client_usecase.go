package usecase

import (
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
)

type ClientUsecase struct {
	ClientRepo repository.ClientRepository
}

func NewClientUseCase(clientRepo repository.ClientRepository) *ClientUsecase {
	return &ClientUsecase{
		ClientRepo: clientRepo,
	}
}

func (uc *ClientUsecase) CreateClient(client *model.Client) (*model.Client, error) {
	return uc.ClientRepo.Create(client)
}

func (uc *ClientUsecase) GetClients() ([]*model.Client, error) {
	return uc.ClientRepo.GetAll()
}
