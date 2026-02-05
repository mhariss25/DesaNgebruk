package controller

import (
	"os"
	"path/filepath"
	"time"

	"DesaNgebruk/database"
	"DesaNgebruk/models"

	"github.com/gofiber/fiber/v2"
)

func UploadImages(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse multipart form", "details": err.Error()})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "images is required"})
	}

	// pastikan folder ada
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create uploads folder", "details": err.Error()})
	}

	base := c.BaseURL()
	urls := make([]string, 0, len(files))

	for _, file := range files {
		uniqueFilename := generateUniqueFilename(file.Filename)
		dst := filepath.Join("uploads", uniqueFilename)

		// simpan file ke disk
		if err := c.SaveFile(file, dst); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save file", "details": err.Error()})
		}

		// simpan info gambar ke MySQL
		img := models.Gambar{
			OriginalName: file.Filename,
			Path:         uniqueFilename, // simpan nama file saja
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := database.DB.Create(&img).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save image to DB", "details": err.Error()})
		}

		// url yang bisa diakses frontend
		urls = append(urls, base+"/uploads/"+uniqueFilename)
	}

	return c.JSON(fiber.Map{
		"message":        "Images uploaded successfully",
		"urls":           urls,
		"presigned_urls": urls, // biar FE lama tetap jalan
	})

}
