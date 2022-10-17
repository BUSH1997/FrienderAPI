package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/vk_api"
	"github.com/sirupsen/logrus"
)

type ImageUseCase struct {
	eventRepository event.Repository
	imageRepository image.Repository
	logger          *logrus.Logger
	vk              vk_api.VKApi
}

func New(imageRepository image.Repository, eventRepository event.Repository, logger *logrus.Logger, vk vk_api.VKApi) image.UseCase {
	return &ImageUseCase{
		eventRepository: eventRepository,
		imageRepository: imageRepository,
		logger:          logger,
		vk:              vk,
	}
}
