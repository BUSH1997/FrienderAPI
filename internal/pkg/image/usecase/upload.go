package usecase

import (
	"context"
	"github.com/labstack/gommon/log"
	"mime/multipart"
)

func (us *ImageUseCase) UploadImage(ctx context.Context, file *multipart.FileHeader, uid string) error {
	err := us.imageRepository.UploadImage(ctx, file)
	if err != nil {
		log.Error(err)
		return err
	}

	stringVkId, err := us.vk.UploadPhoto(file)
	if err != nil {
		log.Error(err)
	}
	link := "https://friender.hb.bizmrg.com/" + file.Filename

	return us.eventRepository.UploadAvatar(ctx, uid, link, stringVkId)
}
