package models

type Kategori struct {
	Id_Kategori   uint   `json:"id_kategori" gorm:"primaryKey"`
	Kategori_name string `json:"kategori_name"`
}
