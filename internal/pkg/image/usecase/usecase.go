package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/sirupsen/logrus"
)

type ImageUseCase struct {
	eventRepository event.Repository
	imageRepository image.Repository
	logger          *logrus.Logger
}

func New(imageRepository image.Repository, eventRepository event.Repository, logger *logrus.Logger) image.UseCase {
	return &ImageUseCase{
		eventRepository: eventRepository,
		imageRepository: imageRepository,
		logger:          logger,
	}
}
