package models

import (
	"time"

	"gorm.io/gorm"
)

type Blogger struct {
	Id_Blogger      uint           `json:"id_blogger" gorm:"primaryKey"`
	Heading_Blogger string         `json:"heading_bloger"`
	KategoriID      uint           `json:"kategori_id" gorm:"column:kategori_id"`
	Kategori        Kategori       `json:"kategori" gorm:"foreignKey:KategoriID"`
	Name_Blog       string         `json:"name_blog" form:"name_blog"`
	FillBlogger     string         `json:"fill_blogger" form:"fill_blogger"`
	Images          []Image        `json:"images" gorm:"foreignKey:BlogID"`
	User_Id         uint           `json:"user_id"`
	User            User           `json:"user" gorm:"foreignKey:User_Id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"delete_at" gorm:" softDelete: true"  format:"date-time"`
}
