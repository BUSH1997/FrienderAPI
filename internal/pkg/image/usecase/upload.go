package usecase

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"mime/multipart"
	"strings"
)

func (us *ImageUseCase) UploadImage(ctx context.Context, files map[string][]*multipart.FileHeader, uid string) error {
	for i := 0; i < len(files); i++ {
		currentFieldName := fmt.Sprintf("photo%d", i)
		err := us.imageRepository.UploadImage(ctx, files[currentFieldName][0])
		if err != nil {
			log.Error(err)
			return err
		}
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

	links := ""
	for i := 1; i < len(files); i++ {
		currentFieldName := fmt.Sprintf("photo%d", i)
		links += "https://friender.hb.bizmrg.com/" + files[currentFieldName][0].Filename + ","
	}

	if links == "" {
		return nil
	}
	links = strings.TrimSuffix(links, ",")
	return us.eventRepository.UploadImage(ctx, uid, links)
}
