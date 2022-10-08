package image

import (
	"context"
	"mime/multipart"
)

type Repository interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) error
}
