package database

import (
	"context"
	"image"
	"watermark-service/internal"
)

type Service interface {
	Add(ctx context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error)
	Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error)
	Remove(ctx context.Context, ticketID string) (int, error)
	ServiceStatus(ctx context.Context) (int, error)
}
