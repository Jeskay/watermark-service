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
	span := internal.StartSpan("picture generation", ctx)
	defer span.Finish()
	if text == "" && logo == nil {
		return nil, errors.New("No data to insert")
	}
	watermark := internal.CombineTextWithLogo(logo, text)
	w.log.Info("Logo creation", zap.String("Status", "Complete"))
	if fill {
		w.log.Info("Fill image", zap.String("Status", "Started"))
		return internal.FillImageWithWatermarks(watermark, Image), nil
	}
	w.log.Info("Add watermark to image", zap.String("Status", "Started"))
	return internal.AddWatermarkToImage(watermark, Image, pos), nil
}

func (w *pictureService) ServiceStatus(ctx context.Context) (int64, error) {
	span := internal.StartSpan("status retrieval", ctx)
	defer span.Finish()
	w.log.Info("Request", zap.Int("Status", http.StatusOK))
	return http.StatusOK, nil
}
