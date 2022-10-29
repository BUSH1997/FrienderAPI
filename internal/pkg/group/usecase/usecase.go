package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/sirupsen/logrus"
)

type groupUseCase struct {
	logger         *logrus.Logger
	repository     group.Repository
	repositoryUser profile.Repository
}

func New(logger *logrus.Logger, repository group.Repository, repositoryProfile profile.Repository) group.UseCase {
	return &groupUseCase{
		logger:         logger,
		repository:     repository,
		repositoryUser: repositoryProfile,
	}
}
