package watermark

import (
	"context"
	"watermark-service/internal"
)

type Service interface {
	Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error)
	Status(ctx context.Context, ticketID int64) (internal.Status, error)
	Watermark(ctx context.Context, ticketID int64, mark string) (int, error)
	AddDocument(ctx context.Context, doc *internal.Document) (int64, error)
	ServiceStatus(ctx context.Context) (int, error)
}
