package picture

import (
	"context"
	"image"
	"watermark-service/internal"
)

type Service interface {
	Create(ctx context.Context, Image image.Image, logo image.Image, text string, fill bool, pos internal.Position) (image.Image, error)
	ServiceStatus(ctx context.Context) (int64, error)
}
