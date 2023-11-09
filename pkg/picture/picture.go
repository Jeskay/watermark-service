package picture

import (
	"context"
	"image"
	"image/draw"
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
	rect := Image.Bounds()
	bg := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	var offset image.Point
	switch pos {
	case internal.LeftTop:
		offset = image.Pt(0, 0)
	case internal.RightTop:
		offset = image.Pt(rect.Dx()-logo.Bounds().Dx(), 0)
	case internal.LeftBottom:
		offset = image.Pt(0, rect.Dy()-logo.Bounds().Dy())
	case internal.RightBottom:
		offset = image.Pt(rect.Dx()-logo.Bounds().Dx(), rect.Dy()-logo.Bounds().Dy())
	default:
		offset = image.Pt(0, 0)
	}
	draw.Draw(bg, Image.Bounds(), Image, image.Point{0, 0}, draw.Over)
	draw.Draw(bg, Image.Bounds().Add(offset), logo, image.Point{0, 0}, draw.Over)

	return bg, nil
}

func (w *pictureService) ServiceStatus(_ context.Context) (int64, error) {
	w.log.Info("Checking the service health...")
	return http.StatusOK, nil
}
