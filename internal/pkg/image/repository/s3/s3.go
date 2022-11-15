package s3

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type ImageRepository struct {
	logger hardlogger.Logger
}

func New(logger hardlogger.Logger) image.Repository {
	return &ImageRepository{
		logger: logger,
	}
}
