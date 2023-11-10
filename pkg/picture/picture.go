package picture

import (
	"context"
	"errors"
	"image"

	"log/slog"
	"net/http"
	"watermark-service/internal"
)

type pictureService struct {
	log *slog.Logger
}

func NewService() Service {
	return &pictureService{
		log: slog.Default(),
	}
}

func (w *pictureService) Create(ctx context.Context, Image image.Image, logo image.Image, text string, fill bool, pos internal.Position) (image.Image, error) {
	if text == "" && logo == nil {
		return nil, errors.New("No data to insert")
	}
	watermark := internal.CombineTextWithLogo(logo, text)
	return internal.AddWatermarkToImage(watermark, Image, pos), nil
}

func (w *pictureService) ServiceStatus(_ context.Context) (int64, error) {
	w.log.Info("Checking the service health...")
	return http.StatusOK, nil
}
