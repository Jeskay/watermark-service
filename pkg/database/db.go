package database

import (
	"context"
	"image"
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

func (d *dbService) Add(_ context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error) {
	return "", nil
}

func (d *dbService) Get(_ context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	var result []internal.Document
	d.orm.Model(&database.Document{}).Find(result)
	return result, nil
}

func (d *dbService) Remove(_ context.Context, ticketId string) (int, error) {
	d.orm.Delete(&database.Document{}, ticketId)
	return http.StatusOK, nil
}

func (d *dbService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
