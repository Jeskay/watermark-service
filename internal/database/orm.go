package database

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	AuthorId uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Title    string    `gorm:"type:varchar(255);not null"`
	ImageUrl string    `gorm:"type:text;not null"`
}

func (d *Document) BeforeCreate(*gorm.DB) error {
	d.ID = uuid.New()

	return nil
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
