package usecase

import (
	"context"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/vk_api"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"mime/multipart"
	"strings"
)

func (uc *ImageUseCase) UploadImage(ctx context.Context, files map[string][]*multipart.FileHeader, uid string) error {
	ctx = uc.logger.WithCaller(ctx)

	for i := 0; i < len(files); i++ {
		currentFieldName := fmt.Sprintf("photo%d", i)
		err := uc.imageRepository.UploadImage(ctx, files[currentFieldName][0])
		if err != nil {
			log.Error(err)
			return err
		}
	}

	stringVkId, err := uc.vk.UploadPhoto(files["photo0"][0], vk_api.UploadPhotoParam{Type: vk_api.Default})
	if err != nil {
		log.Error(err)
	}
	linkAvatar := "https://friender.hb.bizmrg.com/" + files["photo0"][0].Filename
	err = uc.eventRepository.UploadAvatar(ctx, uid, linkAvatar, stringVkId)
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
	return uc.eventRepository.UploadImage(ctx, uid, links)
}

func (uc *ImageUseCase) UploadImageAlbum(ctx context.Context, form *multipart.Form) (interface{}, error) {
	ctx = uc.logger.WithCaller(ctx)

	uploadServer := form.Value["upload_server"]
	if uploadServer == nil {
		uc.logger.WithCtx(ctx).Errorf("Empty upload_server")
		return []string{}, errors.New("Empty upload_server")
	}
	fmt.Println(form)
	photos := form.File

	if photos["photos"] == nil {
		uc.logger.WithCtx(ctx).Errorf("Empty photos")
		return []string{}, errors.New("Empty photos")
	}

	respServer, err := uc.vk.UploadPhotoOnUriServer(photos["photos"][0], uploadServer[0])
	if err != nil {
		uc.logger.WithCtx(ctx).Errorf("Error upload photo album")
		return []string{}, errors.New("Error Upload photo album")
	}

	return respServer, nil
}
