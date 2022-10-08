package image

import (
	"context"
	"mime/multipart"
)

type RepositoryFS interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader) error
}

type RepositoryBD interface {
	UploadImage(ctx context.Context, uid string, fileName string) error
}
