package entity

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        int            `gorm:"primarykey"`
	CreatedAt time.Time      `                  json:"-"`
	UpdatedAt time.Time      `                  json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"      json:"-"`
}
