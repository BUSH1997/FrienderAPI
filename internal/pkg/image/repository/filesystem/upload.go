package filesystem

import (
	"context"
	"io"
	"mime/multipart"
	"os"
)

func (r *ImageRepository) UploadImage(ctx context.Context, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create("/home/ubuntu/testfriender/static/" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
