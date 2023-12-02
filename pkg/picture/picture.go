package picture

import (
	"context"
	"errors"
	"image"

	"net/http"
	"watermark-service/internal"

	"go.uber.org/zap"
)

type pictureService struct {
	log *zap.Logger
}

func NewService() Service {
	return &pictureService{
		log: zap.L().With(zap.String("Service", "PictureService")),
	}
}

func (w *pictureService) Create(ctx context.Context, Image image.Image, logo image.Image, text string, fill bool, pos internal.Position) (image.Image, error) {
	if text == "" && logo == nil {
		return nil, errors.New("No data to insert")
	}
	watermark := internal.CombineTextWithLogo(logo, text)
	if fill {
		return internal.FillImageWithWatermarks(watermark, Image), nil
	}
	return internal.AddWatermarkToImage(watermark, Image, pos), nil
}

func (w *pictureService) ServiceStatus(_ context.Context) (int64, error) {
	w.log.Info("Request", zap.Int("Status", http.StatusOK))
	return http.StatusOK, nil
}
