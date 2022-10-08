package filesystem

import (
	"github.com/sirupsen/logrus"
)

type ImageRepository struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) ImageRepository {
	return ImageRepository{
		logger: logger,
	}
}
