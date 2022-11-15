package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type groupUseCase struct {
	logger         hardlogger.Logger
	repository     group.Repository
	repositoryUser profile.Repository
}

func New(logger hardlogger.Logger, repository group.Repository, repositoryProfile profile.Repository) group.UseCase {
	return &groupUseCase{
		logger:         logger,
		repository:     repository,
		repositoryUser: repositoryProfile,
	}
}
