package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id_User   uint           `json:"id_user" gorm:"primaryKey"`
	Email     string         `json:"email" form:"email"`
	Nama      string         `json:"nama" form:"nama"`
	Password  string         `json:"password" form:"password"`
	Username  string         `json:"username" form:"username"`
	Role      string         `json:"role" form:"role"`
	CreatedAt time.Time      `json:"created_at" form:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"delete_at" gorm:""  format:"date-time"`
}
