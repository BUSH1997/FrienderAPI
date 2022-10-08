package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/sirupsen/logrus"
)

type ImageUseCase struct {
	repositoryFS       image.RepositoryFS
	repositoryPostgres image.RepositoryBD
	logger             *logrus.Logger
}

func New(repositoryFileSystem image.RepositoryFS, repositoryPostgres image.RepositoryBD, logger *logrus.Logger) image.UseCase {
	return &ImageUseCase{
		repositoryFS:       repositoryFileSystem,
		repositoryPostgres: repositoryPostgres,
		logger:             logger,
	}
}
