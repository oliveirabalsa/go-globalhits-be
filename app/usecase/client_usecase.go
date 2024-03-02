package usecase

import (
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
)

type ClientUsecase struct {
	ClientRepo repository.ClientRepository
}

func (uc *ClientUsecase) CreateClient(client *model.Client) (*model.Client, error) {
	return uc.ClientRepo.Create(client)
}
