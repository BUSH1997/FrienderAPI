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

func (uc *ImageUseCase) UploadImageAlbum(ctx context.Context, form *multipart.Form) ([]string, error) {
	ctx = uc.logger.WithCaller(ctx)

	token := form.Value["token"]
	if token == nil {
		uc.logger.WithCtx(ctx).Errorf("Empty user key")
		return []string{}, errors.New("Empty user key")
	}

	albumId := form.Value["album_id"]
	if albumId == nil {
		uc.logger.WithCtx(ctx).Errorf("Empty album_id")
		return []string{}, errors.New("Empty album_id")
	}

	idPhotos := make([]string, 0)
	photos := form.File
	for _, v := range photos["photos"] {
		stringVkId, err := uc.vk.UploadPhoto(v, vk_api.UploadPhotoParam{Type: vk_api.Album, AlbumId: albumId[0], Token: token[0]})
		if err != nil {
			uc.logger.WithCtx(ctx).Errorf("Error upload photo album")
			return []string{}, errors.New("Error Upload photo album")
		}
		idPhotos = append(idPhotos, stringVkId)
	}

	return idPhotos, nil
}
