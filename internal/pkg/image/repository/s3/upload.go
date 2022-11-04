package s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime/multipart"
	"time"
)

func (r *ImageRepository) UploadImage(ctx context.Context, file *multipart.FileHeader, filename string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	b, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	var timeout time.Duration

	bucket := "friender"
	timeout = 60000000000
	endpoint := "hb.bizmrg.com"

	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region:   aws.String("ru-msk"),
				Endpoint: &endpoint,
				Credentials: credentials.NewStaticCredentials("rABxzdZ2uErHGWTxFZ4GPw",
					"7a94iZyVdGABZRT4v395pWpr6iB5fpYMivd2ruPiEFZq", ""),
			},
		))

	svc := s3.New(sess)

	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	defer cancelFn()

	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(b),
	})
	if err, ok := err.(awserr.Error); ok && err.Code() == request.CanceledErrorCode {
		return errors.Wrap(err, "upload canceled due to timeout")
	}
	if err != nil {
		return errors.Wrap(err, "failed to upload object")
	}

	return nil
}
