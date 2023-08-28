package database

import (
	"context"
	"net/http"
	"watermark-service/internal"
	"watermark-service/internal/database"

	"gorm.io/gorm"
)

type dbService struct {
	orm *gorm.DB
}

func NewService(dbORM *gorm.DB) *dbService {
	return &dbService{orm: dbORM}
}

func (d *dbService) Add(_ context.Context, doc *internal.Document) (int64, error) {
	document := database.NewDocument(doc)
	result := d.orm.Model(&database.Document{}).Create(document)
	if result.Error != nil {
		return -1, result.Error
	}
	return document.TicketID, nil
}

func (d *dbService) Get(_ context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	var result []internal.Document
	d.orm.Model(&database.Document{}).Find(result)
	return result, nil
}

func (d *dbService) Update(_ context.Context, ticketId int64, doc *internal.Document) (int, error) {
	document := database.NewDocument(doc)
	document.TicketID = ticketId
	d.orm.Model(&database.Document{}).Save(document)
	return http.StatusOK, nil
}

func (d *dbService) Remove(_ context.Context, ticketId int64) (int, error) {
	d.orm.Delete(&database.Document{}, ticketId)
	return http.StatusOK, nil
}

func (d *dbService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
