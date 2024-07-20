package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"type:varchar(100);not null" json:"title"`
	ISBN          string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"isbn"`
	PublishedDate time.Time      `gorm:"not null" json:"publishedDate"`
	AuthorID      uint           `gorm:"not null" json:"authorID"`
	Author        Author         `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"author"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
