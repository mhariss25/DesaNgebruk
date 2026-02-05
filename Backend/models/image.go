package models

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	OriginalName string         `json:"original_name"`
	Path         string         `json:"images"`
	BlogID       uint           `json:"blog_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"delete_at" gorm:"softDelete:true"  format:"date-time"`
}
