package watermark

import (
	"context"
	"errors"
	"io"
	"regexp"
	"watermark-service/internal"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Storage interface {
	Upload(ctx context.Context, name string, image io.Reader) (string, error)
	Delete(ctx context.Context, id string) error
}

type CloudinaryStorage struct {
	instance *cloudinary.Cloudinary
}

func NewCloudinaryStorage(cloud, apiKey, secretKey string) *CloudinaryStorage {
	cld, err := cloudinary.NewFromParams(cloud, apiKey, secretKey)
	if err != nil {
		panic(err)
	}
	return &CloudinaryStorage{instance: cld}
}

func (s *CloudinaryStorage) Upload(ctx context.Context, name string, image io.Reader) (string, error) {
	span := internal.StartSpan("Cloudinary upload", ctx)
	defer span.Finish()
	res, err := s.instance.Upload.Upload(ctx, image, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return res.URL, err
}

func (s *CloudinaryStorage) Delete(ctx context.Context, url string) error {
	id, ok := idFromURL(url)
	if !ok {
		return errors.New("Invalid url")
	}
	_, err := s.instance.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: id})
	if err != nil {
		return err
	}
	return nil
}

func idFromURL(url string) (string, bool) {
	r, err := regexp.Compile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	if err != nil {
		return "", false
	}
	matches := r.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[2], true
	}
	return "", false

}
