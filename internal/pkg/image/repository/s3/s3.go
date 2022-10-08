package s3

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/sirupsen/logrus"
)

type ImageRepository struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) image.Repository {
	return &ImageRepository{
		logger: logger,
	}
}
