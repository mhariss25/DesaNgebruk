package migrations

import (
	"fmt"
	"log"

	"DesaNgebruk/database"
	"DesaNgebruk/models"
)

func MigrationTable() {

	err := database.DB.AutoMigrate(&models.User{}, &models.Blogger{}, &models.Image{}, &models.Kategori{}, &models.Gambar{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Migrated successfully")

}
