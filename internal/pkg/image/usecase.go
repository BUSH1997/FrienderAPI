package image

import (
	"context"
	"mime/multipart"
)

type UseCase interface {
	UploadImage(ctx context.Context, file map[string][]*multipart.FileHeader, uid string) error
	UploadImageAlbum(ctx context.Context, form *multipart.Form) (interface{}, error)
}
