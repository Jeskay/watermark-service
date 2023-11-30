package watermark

import (
	"github.com/google/uuid"
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

func InitDb(db *gorm.DB) error {
	return db.AutoMigrate(&Document{})
}
