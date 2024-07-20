package models

import (
	"gorm.io/gorm"
)

type Author struct {
	// The primary key is of an integer type so GORM will automatically set it to auto-increment (if i didn't enter)
	// If i didn't write json it will define the property during Serialization and Deserialization as it is (capitalized)
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
