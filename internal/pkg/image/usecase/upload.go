package usecase

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"mime/multipart"
	"strings"
)

func (us *ImageUseCase) UploadImage(ctx context.Context, files map[string][]*multipart.FileHeader, uid string) error {
	var links string
	for i := 0; i < len(files); i++ {
		currentFieldName := fmt.Sprintf("photo%d", i)
		fileName, err := uuid.NewV4()
		if err != nil {
			return errors.Wrap(err, "failed to generate filename")
		}

		err = us.imageRepository.UploadImage(ctx, files[currentFieldName][0], fileName.String())
		if err != nil {
			return err
		}

		links += "https://friender.hb.bizmrg.com/" + fileName.String() + ","
	}

	stringVkId, err := us.vk.UploadPhoto(files["photo0"][0])
	if err != nil {
		log.Error(err)
	}
	linkAvatar := "https://friender.hb.bizmrg.com/" + files["photo0"][0].Filename
	err = us.eventRepository.UploadAvatar(ctx, uid, linkAvatar, stringVkId)
	if err != nil {
		log.Error(err)
		return err
	}

	if links == "" {
		return nil
	}
	links = strings.TrimSuffix(links, ",")
	return us.eventRepository.UploadImage(ctx, uid, links)
}
