package usecase

import "github.com/BUSH1997/FrienderAPI/internal/pkg/models"

func (uc *UseCase) GetOneProfile(id string) (models.Profile, error) {
	return models.Profile{}, nil
}

func (uc *UseCase) GetAllStatusesUser(id string) ([]models.Status, error) {
	return nil, nil
}
