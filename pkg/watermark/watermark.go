package watermark

import (
	"context"
	"log/slog"
	"net/http"
	"watermark-service/internal"
)

type watermarkService struct {
	log *slog.Logger
}

func NewService() Service {
	return &watermarkService{log: slog.Default()}
}

func (w *watermarkService) Get(_ context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	// there must be real query from database
	doc := internal.Document{
		Content: "book",
		Title:   "Security Analysis",
		Author:  "B. Graham",
		Topic:   "Finance",
	}
	return []internal.Document{doc}, nil
}

func (w *watermarkService) Status(_ context.Context, ticketID string) (internal.Status, error) {
	//there must be real query from database with document info
	return internal.InProgress, nil
}

func (w *watermarkService) Watermark(_ context.Context, ticketID string, mark string) (int, error) {
	//update db with watermark field as non empty
	return http.StatusOK, nil
}

func (w *watermarkService) AddDocument(_ context.Context, doc *internal.Document) (string, error) {
	newTicketID := "12321"
	return newTicketID, nil
}

func (w *watermarkService) ServiceStatus(_ context.Context) (int, error) {
	w.log.Info("Checking the service health...")
	return http.StatusOK, nil
}
