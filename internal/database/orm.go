package database

import (
	"fmt"
	"watermark-service/internal"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	TicketID  int64  `gorm:"primaryKey;autoIncrement"`
	Content   string `gorm:"type:varchar(100)"`
	Title     string `gorm:"type:varchar(100)"`
	Author    string `gorm:"type:varchar(100)"`
	Topic     string `gorm:"type:varchar(100)"`
	Watermark string `gorm:"type:varchar(100)"`
}

func NewDocument(doc *internal.Document) *Document {
	return &Document{
		Content:   doc.Content,
		Title:     doc.Title,
		Author:    doc.Author,
		Topic:     doc.Topic,
		Watermark: doc.Watermark,
	}
}

func (d *Document) ToJSON() *internal.Document {
	return &internal.Document{
		Content:   d.Content,
		Title:     d.Title,
		Author:    d.Author,
		Topic:     d.Topic,
		Watermark: d.Watermark,
	}
}

func Init(host, port, user, dbname, pass string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, pass)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return db, err
	}
	err = db.AutoMigrate(&Document{})
	return db, err
}
