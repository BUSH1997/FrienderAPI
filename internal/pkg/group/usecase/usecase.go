package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/sirupsen/logrus"
)

type groupUseCase struct {
	logger     *logrus.Logger
	repository group.Repository
}

func New(logger *logrus.Logger, repository group.Repository) group.UseCase {
	return &groupUseCase{
		logger:     logger,
		repository: repository,
	}
}
