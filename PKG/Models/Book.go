package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	// uint ==> unsigned integer, It can store positive values and zero.
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"type:varchar(100);not null" json:"title"`
	ISBN          string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"isbn"`
	PublishedDate time.Time      `gorm:"not null" json:"publishedDate"`
	AuthorID      uint           `gorm:"not null" json:"authorID"`
	Author        Author         `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;" json:"author"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	// for handling soft deletes
	// When a record is "deleted," the current timestamp is set in the Time field, and the Valid field is set to true. Records with a Valid value of false are considered not deleted.
	// json:"-"` ==> to ignore this field when marshalling or unmarshalling JSON.
	// Marshaling ==> Converting Go data structures into a format that can be easily stored or transmitted, such as JSON.
	// Unmarshaling ==> Converting data from these formats back into Go data structures.
	// gorm:"index" applied to DeletedAt means that GORM will create an index on the DeletedAt column in the database.
	// Indexes speed up queries that filter or sort by the indexed column.
	// CREATE INDEX idx_deleted_at ON books (deleted_at);
}
