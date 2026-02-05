package controller

import (
	"strconv"

	"github.com/agnarbriantama/DesaNgembruk-Backend/database"
	"github.com/agnarbriantama/DesaNgembruk-Backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllKategori(c *fiber.Ctx) error {
	var kategoris []models.Kategori
	result := database.DB.Find(&kategoris)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(kategoris)
}

// CreateKategori creates a new Kategori entry and saves it to the database
func CreateKategori(c *fiber.Ctx) error {
	// Parse JSON request body into a Kategori struct
	var kategori models.Kategori
	if err := c.BodyParser(&kategori); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Create a new Kategori entry in the database
	result := database.DB.Create(&kategori)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create Kategori", "details": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(kategori)
}

// UpdateKategori updates an existing Kategori entry in the database
func UpdateKategori(c *fiber.Ctx) error {
	kategoriID := c.Params("id")
	// Parse JSON request body into a Kategori struct
	var kategori models.Kategori
	if err := c.BodyParser(&kategori); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Get the existing Kategori by ID from the database
	existingKategori := models.Kategori{}
	result := database.DB.Where("id_kategori = ?", kategoriID).First(&existingKategori)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kategori not found"})
	}

	// Update the existing Kategori with new data
	result = database.DB.Model(&existingKategori).Updates(&kategori)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Kategori", "details": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(existingKategori)
}

// DeleteKategori deletes an existing Kategori entry from the database
func DeleteKategori(c *fiber.Ctx) error {
	kategoriID := c.Params("id")

	// Get the existing Kategori by ID from the database
	existingKategori := models.Kategori{}
	result := database.DB.Where("id_kategori = ?", kategoriID).First(&existingKategori)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kategori not found"})
	}

	// Delete the Kategori from the database
	result = database.DB.Unscoped().Delete(&existingKategori)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete Kategori", "details": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kategori deleted successfully"})
}

func GetCategoryById(c *fiber.Ctx) error {
	// Extract ID from URL parameter
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// If the ID is not a valid integer
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var kategori models.Kategori
	result := database.DB.First(&kategori, id)

	if result.Error != nil {
		// If no record is found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kategori not found"})
	}

	// Return the found category
	return c.JSON(kategori)
}

func BloggerByCategoryHandler(c *fiber.Ctx) error {
	// Ambil id_kategori dari parameter URL
	idKategori := c.Params("id_kategori")

	// Konversi id_kategori ke tipe data uint
	idKategoriUint, err := strconv.ParseUint(idKategori, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}
	kategoriID := uint(idKategoriUint)

	// Fetch bloggers berdasarkan kategori
	var bloggers []models.Blogger
	if err := database.DB.Where("kategori_id = ?", kategoriID).Preload("Kategori").Preload("User").Find(&bloggers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bloggers by category", "details": err.Error()})
	}

	// ... (lakukan operasi lain jika diperlukan)

	response := fiber.Map{
		"message":  "Bloggers retrieved successfully by category",
		"bloggers": bloggers,
	}

	return c.JSON(response)
}
