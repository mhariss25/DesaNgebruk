package controller

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"math"

	"DesaNgebruk/database"
	"DesaNgebruk/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
)

func saveToLocal(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	// pastikan folder ada
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		return "", err
	}

	unique := generateUniqueFilename(file.Filename)
	dst := filepath.Join("uploads", unique)

	if err := c.SaveFile(file, dst); err != nil {
		return "", err
	}

	return unique, nil // simpan ke DB
}

func deleteLocalFile(filename string) error {
	if filename == "" {
		return nil
	}
	path := filepath.Join("uploads", filename)
	// kalau file tidak ada, anggap sukses (biar edit/delete tidak gagal)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(path)
}

func CreateBlogger(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if len(tokenString) < 7 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}
	jwtToken := strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	userID := uint(claims["id_user"].(float64))
	username := claims["username"].(string)

	var blogger models.Blogger
	if err := c.BodyParser(&blogger); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	blogger.User_Id = userID

	kategoriID := c.FormValue("kategori_id")
	if kategoriID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kategori ID is required"})
	}

	kategoriIDInt, err := strconv.Atoi(kategoriID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Kategori ID"})
	}

	var kategori models.Kategori
	if err := database.DB.First(&kategori, kategoriIDInt).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Kategori ID"})
	}
	blogger.KategoriID = uint(kategoriIDInt)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}

	// heading
	headingFiles := form.File["heading_blogger"]
	if len(headingFiles) > 0 {
		headingFile := headingFiles[0]
		if !isValidImageFormat(headingFile) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid heading image format"})
		}

		filename, err := saveToLocal(c, headingFile)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save heading locally", "details": err.Error()})
		}
		blogger.Heading_Blogger = filename
	}

	files := form.File["images"]

	// kalau ada images, pakai transaksi biar konsisten
	if len(files) > 0 {
		tx := database.DB.Begin()

		if err := tx.Create(&blogger).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create blogger", "details": err.Error()})
		}

		for _, file := range files {
			if !isValidImageFormat(file) {
				tx.Rollback()
				return c.Status(400).JSON(fiber.Map{"error": "Invalid image format"})
			}

			filename, err := saveToLocal(c, file)
			if err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to save image locally", "details": err.Error()})
			}

			image := models.Image{
				OriginalName: file.Filename,
				Path:         filename,
				BlogID:       blogger.Id_Blogger,
			}

			if err := tx.Create(&image).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to save image to the database", "details": err.Error()})
			}
		}

		tx.Commit()
	} else {
		// tanpa images
		if err := database.DB.Create(&blogger).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create blogger", "details": err.Error()})
		}
	}

	if err := database.DB.Model(&blogger).Association("Kategori").Find(&blogger.Kategori); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to load category data", "details": err.Error()})
	}

	// opsional: balikkan URL siap dipakai FE
	base := c.BaseURL()
	if blogger.Heading_Blogger != "" {
		blogger.Heading_Blogger = base + "/uploads/" + blogger.Heading_Blogger
	}
	if err := database.DB.Model(&blogger).Association("Images").Find(&blogger.Images); err == nil {
		for i := range blogger.Images {
			if blogger.Images[i].Path != "" {
				blogger.Images[i].Path = base + "/uploads/" + blogger.Images[i].Path
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":  "Blogger created successfully",
		"user_id":  userID,
		"username": username,
		"blogger":  blogger,
	})
}

func EditBlogger(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if len(tokenString) < 7 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}
	jwtToken := strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	userID := uint(claims["id_user"].(float64))
	role := claims["role"].(string)
	username := claims["username"].(string)

	var updatedBlogger models.Blogger
	if err := c.BodyParser(&updatedBlogger); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	kategoriID := c.FormValue("kategori_id")
	if kategoriID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kategori ID is required"})
	}
	kategoriIDInt, err := strconv.Atoi(kategoriID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Kategori ID"})
	}

	var kategori models.Kategori
	if err := database.DB.First(&kategori, kategoriIDInt).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Kategori ID"})
	}

	var existingBlogger models.Blogger
	if err := database.DB.Preload("Images").Where("id_blogger = ?", c.Params("id_blogger")).First(&existingBlogger).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blogger not found"})
	}

	if !(role == "admin" || (role == "writter" && existingBlogger.User_Id == userID)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}

	// heading baru (kalau ada) -> hapus heading lama dulu
	headingFiles := form.File["heading_blogger"]
	if len(headingFiles) > 0 {
		headingFile := headingFiles[0]
		if !isValidImageFormat(headingFile) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid heading image format"})
		}

		// hapus heading lama (kalau ada dan masih filename lokal)
		_ = deleteLocalFile(existingBlogger.Heading_Blogger)

		filename, err := saveToLocal(c, headingFile)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save heading locally", "details": err.Error()})
		}
		existingBlogger.Heading_Blogger = filename
	}

	// hapus semua images lama (file + row DB)
	for _, oldImage := range existingBlogger.Images {
		_ = deleteLocalFile(oldImage.Path)

		if err := database.DB.Delete(&oldImage).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete old image from database", "details": err.Error()})
		}
	}
	existingBlogger.Images = []models.Image{}

	// upload images baru
	files := form.File["images"]
	for _, file := range files {
		if !isValidImageFormat(file) {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid image format"})
		}

		// Cek duplikasi (opsional, tapi sesuai kode kamu)
		var duplicateImage models.Image
		result := database.DB.Where("original_name = ? AND blog_id = ?", file.Filename, existingBlogger.Id_Blogger).First(&duplicateImage)
		if result.Error == nil {
			// sudah ada -> skip
			continue
		}

		filename, err := saveToLocal(c, file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save image locally", "details": err.Error()})
		}

		newImage := models.Image{
			OriginalName: file.Filename,
			Path:         filename,
			BlogID:       existingBlogger.Id_Blogger,
		}

		if err := database.DB.Create(&newImage).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save image to the database", "details": err.Error()})
		}
		existingBlogger.Images = append(existingBlogger.Images, newImage)
	}

	existingBlogger.Name_Blog = updatedBlogger.Name_Blog
	existingBlogger.FillBlogger = updatedBlogger.FillBlogger
	existingBlogger.KategoriID = uint(kategoriIDInt)

	if err := database.DB.Model(&existingBlogger).Omit("CreatedAt").Save(&existingBlogger).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update blogger", "details": err.Error()})
	}

	// balikkan URL siap pakai FE
	base := c.BaseURL()
	if existingBlogger.Heading_Blogger != "" {
		existingBlogger.Heading_Blogger = base + "/uploads/" + existingBlogger.Heading_Blogger
	}
	for i := range existingBlogger.Images {
		if existingBlogger.Images[i].Path != "" {
			existingBlogger.Images[i].Path = base + "/uploads/" + existingBlogger.Images[i].Path
		}
	}

	return c.JSON(fiber.Map{
		"message":  "Blogger updated successfully",
		"user_id":  userID,
		"username": username,
		"blogger":  existingBlogger,
	})
}

func GetAllBloggers(c *fiber.Ctx) error {
	var bloggers []models.Blogger
	var total int64

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid page number"})
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid page size"})
	}

	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	offset := (page - 1) * pageSize

	searchQuery := c.Query("search", "")
	selectedCategory := c.Query("category", "")

	queryBuilder := database.DB.Model(&models.Blogger{}).
		Joins("left join kategoris on kategoris.id_kategori = bloggers.kategori_id").
		Joins("left join users on users.id_user = bloggers.user_id")

	if searchQuery != "" {
		queryBuilder = queryBuilder.Where(
			"bloggers.name_blog LIKE ? OR kategoris.kategori_name LIKE ? OR users.nama LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%",
		)
	}

	if startDate != "" && endDate != "" {
		endDateWithTime := endDate + " 23:59:59"
		queryBuilder = queryBuilder.Where("bloggers.created_at BETWEEN ? AND ?", startDate, endDateWithTime)
	}

	if selectedCategory != "" {
		queryBuilder = queryBuilder.Where("kategoris.id_kategori = ?", selectedCategory)
	}

	queryBuilder.Count(&total)

	if err := queryBuilder.Preload("Images").Preload("Kategori").Preload("User").
		Offset(offset).Limit(pageSize).Find(&bloggers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve bloggers", "details": err.Error()})
	}

	base := c.BaseURL()
	for i := range bloggers {
		if bloggers[i].Heading_Blogger != "" {
			bloggers[i].Heading_Blogger = base + "/uploads/" + bloggers[i].Heading_Blogger
		}
		for j := range bloggers[i].Images {
			if bloggers[i].Images[j].Path != "" {
				bloggers[i].Images[j].Path = base + "/uploads/" + bloggers[i].Images[j].Path
			}
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(fiber.Map{
		"message":     "All bloggers retrieved successfully",
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalItems":  total,
		"bloggers":    bloggers,
	})
}

func GetBloggerByID(c *fiber.Ctx) error {
	id := c.Params("id_blogger")
	var blogger models.Blogger

	if err := database.DB.Preload("Images").Preload("Kategori").Preload("User").First(&blogger, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blogger not found"})
	}

	base := c.BaseURL()
	if blogger.Heading_Blogger != "" {
		blogger.Heading_Blogger = base + "/uploads/" + blogger.Heading_Blogger
	}
	for i := range blogger.Images {
		if blogger.Images[i].Path != "" {
			blogger.Images[i].Path = base + "/uploads/" + blogger.Images[i].Path
		}
	}

	return c.JSON(blogger)
}

func findImageByPath(path string) (*models.Image, error) {
	var image models.Image
	if err := database.DB.Where("path = ?", path).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func getBloggersByUserID(c *fiber.Ctx, userID int) error {
	var bloggers []models.Blogger

	if err := database.DB.Where("user_id = ?", userID).
		Preload("Images").Preload("Kategori").Preload("User").
		Find(&bloggers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve bloggers for the specified user", "details": err.Error()})
	}

	base := c.BaseURL()
	for i := range bloggers {
		if bloggers[i].Heading_Blogger != "" {
			bloggers[i].Heading_Blogger = base + "/uploads/" + bloggers[i].Heading_Blogger
		}
		for j := range bloggers[i].Images {
			if bloggers[i].Images[j].Path != "" {
				bloggers[i].Images[j].Path = base + "/uploads/" + bloggers[i].Images[j].Path
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":  "Bloggers for the specified user retrieved successfully",
		"bloggers": bloggers,
	})
}

func UserBloggersHandler(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"error": "No authorization token provided"})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	var userID int
	if id, ok := claims["id_user"].(float64); ok {
		userID = int(id)
	} else if id, ok := claims["id_user"].(string); ok {
		userID, err = strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "User ID in token is not an integer"})
		}
	} else {
		return c.Status(401).JSON(fiber.Map{"error": "User ID not found in token"})
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid page number"})
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "5"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid page size"})
	}

	offset := (page - 1) * pageSize

	bloggers, total, err := getPaginatedBloggersByUserID(userID, offset, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve bloggers for the specified user", "details": err.Error()})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	base := c.BaseURL()
	for i := range bloggers {
		if bloggers[i].Heading_Blogger != "" {
			bloggers[i].Heading_Blogger = base + "/uploads/" + bloggers[i].Heading_Blogger
		}
		for j := range bloggers[i].Images {
			if bloggers[i].Images[j].Path != "" {
				bloggers[i].Images[j].Path = base + "/uploads/" + bloggers[i].Images[j].Path
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":     "Bloggers for the specified user retrieved successfully",
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalItems":  total,
		"bloggers":    bloggers,
	})
}

func getPaginatedBloggersByUserID(userID, offset, pageSize int) ([]models.Blogger, int, error) {
	var bloggers []models.Blogger
	var total int64

	if err := database.DB.Where("user_id = ?", userID).
		Preload("Images").Preload("Kategori").Preload("User").
		Offset(offset).Limit(pageSize).
		Find(&bloggers).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Model(&models.Blogger{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return bloggers, int(total), nil
}

func DeleteBlogger(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if len(tokenString) < 7 {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}
	jwtToken := strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	bloggerID := c.Params("id_blogger")

	userID := uint(claims["id_user"].(float64))
	role := claims["role"].(string)

	var blogger models.Blogger
	if err := database.DB.Where("id_blogger = ?", bloggerID).First(&blogger).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blogger not found"})
	}

	if !(role == "admin" || (role == "writter" && blogger.User_Id == userID)) {
		return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
	}

	// hapus heading lokal (kalau ada)
	_ = deleteLocalFile(blogger.Heading_Blogger)

	// hapus images lokal (kalau ada)
	var images []models.Image
	_ = database.DB.Where("blog_id = ?", blogger.Id_Blogger).Find(&images).Error
	for _, img := range images {
		_ = deleteLocalFile(img.Path)
		_ = database.DB.Delete(&img).Error
	}

	if err := database.DB.Delete(&blogger).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete blogger from database", "details": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Blogger deleted successfully",
		"user_id": userID,
		"blogger": blogger,
	})
}

// CustomClaims holds custom claims for JWT.
type CustomClaims struct {
	UserID uint `json:"id_user"`
	jwt.StandardClaims
}

func extractUserClaimsFromToken(c *fiber.Ctx) (*jwt.Token, *CustomClaims, error) {
	tokenString := c.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil
	})
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("Invalid token claims")
	}

	return token, claims, nil
}

func isValidImageFormat(file *multipart.FileHeader) bool {
	allowedFormats := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/jpg":  true,
	}

	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return false
	}

	detectedMIME := strings.ToLower(mime.String())
	return allowedFormats[detectedMIME]
}

func generateUniqueFilename(originalFilename string) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	fileExt := filepath.Ext(originalFilename)
	return fmt.Sprintf("%d%s", timestamp, fileExt)
}
