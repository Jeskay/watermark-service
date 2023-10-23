package watermark

import (
	"context"
	"image"
	"watermark-service/internal"
)

type Middleware func(Service) Service

func WatermarkMiddleware() Middleware {
	return func(next Service) Service {
		return &watermarkMiddleware{
			next: next,
		}
	}
}

type watermarkMiddleware struct {
	next Service
}

func (m *watermarkMiddleware) Create(ctx context.Context, Image, Logo image.Image, text string, fill bool, pos internal.Position) (image.Image, error) {
	return m.next.Create(ctx, Image, Logo, text, fill, pos)
}

func (m *watermarkMiddleware) ServiceStatus(ctx context.Context) (int64, error) {
	return m.next.ServiceStatus(ctx)
}
