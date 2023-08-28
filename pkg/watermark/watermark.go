package watermark

import (
	"context"
	"log/slog"
	"net/http"
	dbproto "watermark-service/api/v1/protos/db"
	"watermark-service/internal"
	"watermark-service/internal/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type watermarkService struct {
	log         *slog.Logger
	dbAvailable bool
	dbClient    dbproto.DatabaseClient
}

func NewService(dbServiceAddr string) Service {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(dbServiceAddr, opts...)
	if err != nil {
		return &watermarkService{
			log:         slog.Default(),
			dbAvailable: false,
		}
	}
	c := dbproto.NewDatabaseClient(conn)
	return &watermarkService{
		log:         slog.Default(),
		dbClient:    c,
		dbAvailable: true,
	}
}

func (w *watermarkService) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	// there must be real query from database
	if !w.dbAvailable {
		return nil, util.ErrDatabaseServiceUnavailable
	}
	var reqFilters []*dbproto.GetRequest_Filters
	for _, filter := range filters {
		reqFilters = append(reqFilters, &dbproto.GetRequest_Filters{
			Key:   filter.Key,
			Value: filter.Value,
		})
	}
	resp, err := w.dbClient.Get(ctx, &dbproto.GetRequest{
		Filters: reqFilters,
	})
	if err != nil {
		return nil, err
	}
	var docs []internal.Document
	for _, doc := range resp.GetDocuments() {
		docs = append(docs, internal.Document{
			Content:   doc.GetContent(),
			Title:     doc.GetTitle(),
			Author:    doc.GetAuthor(),
			Topic:     doc.GetTopic(),
			Watermark: doc.GetWatermark(),
		})
	}
	return docs, nil
}

func (w *watermarkService) Status(ctx context.Context, ticketID int64) (internal.Status, error) {
	return internal.InProgress, nil
}

func (w *watermarkService) Watermark(ctx context.Context, ticketID int64, mark string) (int, error) {
	if !w.dbAvailable {
		return http.StatusFailedDependency, util.ErrDatabaseServiceUnavailable
	}
	resp, err := w.dbClient.Update(ctx, &dbproto.UpdateRequest{
		TicketID: ticketID,
		Document: &dbproto.Document{
			Watermark: mark,
		},
	})
	return int(resp.GetCode()), err
}

func (w *watermarkService) AddDocument(ctx context.Context, doc *internal.Document) (int64, error) {
	if !w.dbAvailable {
		return -1, util.ErrDatabaseServiceUnavailable
	}
	resp, err := w.dbClient.Add(ctx, &dbproto.AddRequest{
		Document: &dbproto.Document{
			Content:   doc.Content,
			Author:    doc.Author,
			Title:     doc.Title,
			Topic:     doc.Topic,
			Watermark: doc.Watermark,
		},
	})
	if err != nil {
		return -1, err
	}
	return resp.TicketID, nil
}

func (w *watermarkService) ServiceStatus(_ context.Context) (int, error) {
	w.log.Info("Checking the service health...")
	return http.StatusOK, nil
}
