package usecase

import (
	"context"
	"github.com/labstack/gommon/log"
	"mime/multipart"
)

func (us *ImageUseCase) UploadImage(ctx context.Context, file *multipart.FileHeader, uid string) error {
	err := us.repositoryFS.UploadImage(ctx, file)
	if err != nil {
		log.Error(err)
		return err
	}

	return us.repositoryPostgres.UploadImage(ctx, uid, file.Filename)
}
