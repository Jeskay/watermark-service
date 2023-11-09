package picture

import (
	"context"
	"image"
	"watermark-service/internal"
)

type Middleware func(Service) Service

func PictureMiddleware() Middleware {
	return func(next Service) Service {
		return &pictureMiddleware{
			next: next,
		}
	}
}

type pictureMiddleware struct {
	next Service
}

func (m *pictureMiddleware) Create(ctx context.Context, Image, Logo image.Image, text string, fill bool, pos internal.Position) (image.Image, error) {
	return m.next.Create(ctx, Image, Logo, text, fill, pos)
}

func (m *pictureMiddleware) ServiceStatus(ctx context.Context) (int64, error) {
	return m.next.ServiceStatus(ctx)
}
