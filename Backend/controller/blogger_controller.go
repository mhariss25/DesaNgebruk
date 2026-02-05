package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/agnarbriantama/DesaNgembruk-Backend/database"
	"github.com/agnarbriantama/DesaNgembruk-Backend/models"
	"github.com/agnarbriantama/DesaNgembruk-Backend/utils/awsutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dgrijalva/jwt-go"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
)

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

	// Pastikan kategoriID tidak kosong
	if kategoriID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kategori ID is required"})
	}

	// Konversi kategoriID ke tipe data yang sesuai (mungkin uint atau int, sesuai dengan tipe data di database)
	kategoriIDInt, err := strconv.Atoi(kategoriID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Kategori ID"})
	}

	// Lakukan validasi apakah ID kategori yang dipilih valid
	// Jika validasi gagal, Anda dapat mengembalikan respons error
	// Contoh: Memeriksa apakah kategori dengan ID yang diberikan ada di database
	var kategori models.Kategori
	if err := database.DB.First(&kategori, kategoriIDInt).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Kategori ID"})
	}

	// Set ID kategori pada blogger
	blogger.KategoriID = uint(kategoriIDInt)

	// Get the S3 uploader from awsutils package
	_, uploader := awsutils.NewS3ClientUploader()

	// Use the fileupload middleware to handle image uploads
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}

	// Handle heading image upload
	headingFiles := form.File["heading_blogger"]
	if len(headingFiles) > 0 {
		headingFile := headingFiles[0] // Assuming only one heading image
		if !isValidImageFormat(headingFile) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid heading image format"})
		}

		fileContent, err := headingFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open heading file", "details": err.Error()})
		}
		defer fileContent.Close()

		uniqueFilename := generateUniqueFilename(headingFile.Filename)
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("desa-ngebruk"),
			Key:    aws.String(uniqueFilename),
			Body:   fileContent,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload heading to Amazon S3", "details": err.Error()})
		}

		// Save heading path to blogger
		blogger.Heading_Blogger = uniqueFilename
	}

	// Retrieve the uploaded files
	files := form.File["images"]

	// Check if any files were uploaded
	if len(files) > 0 {
		_, uploader := awsutils.NewS3ClientUploader()

		tx := database.DB.Begin()

		if err := tx.Create(&blogger).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create blogger", "details": err.Error()})
		}

		for _, file := range files {
			if !isValidImageFormat(file) {
				fmt.Println("Invalid image format:", file.Header.Get("Content-Type"))
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid image format"})
			}

			// Open the file
			fileContent, err := file.Open()
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file", "details": err.Error()})
			}
			defer fileContent.Close()

			uniqueFilename := generateUniqueFilename(file.Filename)

			// Upload file to Amazon S3
			_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
				Bucket: aws.String("desa-ngebruk"),
				Key:    aws.String(uniqueFilename),
				Body:   fileContent,
			})

			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to Amazon S3", "details": err.Error()})
			}

			// Perubahan di sini: Menghapus "s3://desa-ngebruk" dari savePath
			savePath := uniqueFilename

			image := models.Image{
				OriginalName: file.Filename,
				Path:         savePath,
				BlogID:       blogger.Id_Blogger,
			}

			if err := tx.Create(&image).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image to the database", "details": err.Error()})
			}

			blogger.Images = append(blogger.Images, image)

		}

		tx.Commit()
		if len(files) > 0 {
			presignedURLs := make([]string, 0)

			for _, file := range files {
				uniqueFilename := generateUniqueFilename(file.Filename)

				presignedURL, err := awsutils.GeneratePresignURL("desa-ngebruk", uniqueFilename) // Ganti dengan nama bucket yang sesuai
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presigned URL for image", "details": err.Error()})
				}

				// Include presigned URL in response if defined
				if presignedURL != "" {
					presignedURLs = append(presignedURLs, presignedURL)
				}
			}

		}
	} else {
		if err := database.DB.Create(&blogger).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create blogger", "details": err.Error()})
		}
	}
	if err := database.DB.Model(&blogger).Association("Kategori").Find(&blogger.Kategori); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load category data", "details": err.Error()})
	}

	response := fiber.Map{
		"message":  "Blogger created successfully",
		"user_id":  userID,
		"username": username,
		"blogger":  blogger,
	}
	// Check if image presigned URLs are defined before including them in response
	if imagePresignedURLs, ok := response["image_presigned_urls"].([]string); ok && len(imagePresignedURLs) > 0 {
		response["image_presigned_urls"] = imagePresignedURLs
	}

	return c.JSON(response)
}

func EditBlogger(c *fiber.Ctx) error {
	// Parsing dan validasi token JWT
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

	// Mengambil informasi pengguna dari token JWT
	userID := uint(claims["id_user"].(float64))
	role := claims["role"].(string)
	username := claims["username"].(string)

	// Membaca data blogger dari body request
	var updatedBlogger models.Blogger
	if err := c.BodyParser(&updatedBlogger); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// Parse kategori_id from the request
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

	// Mengambil kategori baru dari permintaan
	//     newCategoryID := updatedBlogger.KategoriID
	//  // Validasi kategori baru
	//     var newCategory models.Kategori
	//     if err := database.DB.Where("id_kategori = ?", newCategoryID).First(&newCategory).Error; err != nil {
	//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid new category"})
	//     }

	// Mencari blogger yang akan diperbarui berdasarkan ID
	var existingBlogger models.Blogger
	if err := database.DB.Preload("Images").Where("id_blogger = ?", c.Params("id_blogger")).First(&existingBlogger).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blogger not found"})
	}

	// Validasi peran pengguna
	if !(role == "admin" || (role == "writter" && existingBlogger.User_Id == userID)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	// Mengganti kategori lama dengan kategori baru
	existingBlogger.KategoriID = uint(kategoriIDInt)

	// Mengelola unggahan gambar
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}
	//images heading
	headingFiles := form.File["heading_blogger"]
	if len(headingFiles) > 0 {
		headingFile := headingFiles[0]
		if !isValidImageFormat(headingFile) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid heading image format"})
		}

		uniqueFilename := generateUniqueFilename(headingFile.Filename)
		s3Path, err := uploadImageToS3("desa-ngebruk", uniqueFilename, headingFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload heading to Amazon S3", "details": err.Error()})
		}
		existingBlogger.Heading_Blogger = s3Path
	}

	// Hapus semua gambar lama terkait blogger dari S3 dan database
	for _, oldImage := range existingBlogger.Images {
		if err := deleteImageFromS3("desa-ngebruk", oldImage.Path); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete old image from AWS S3", "details": err.Error()})
		}

		if err := database.DB.Delete(&oldImage).Error; err !=
			nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete old image from database", "details": err.Error()})
		}
	}

	// Hapus referensi gambar lama dari entitas blogger
	existingBlogger.Images = []models.Image{}
	//images blogger
	files := form.File["images"]
	for _, file := range files {
		// Validasi format file
		if !isValidImageFormat(file) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid image format"})
		}

		// Cek duplikasi gambar berdasarkan OriginalName dan BlogID
		var duplicateImage models.Image
		result := database.DB.Where("original_name = ? AND blog_id = ?", file.Filename, existingBlogger.Id_Blogger).First(&duplicateImage)

		if result.Error != nil {
			// Tidak ada gambar duplikat, lakukan proses upload
			uniqueFilename := generateUniqueFilename(file.Filename)
			s3Path, err := uploadImageToS3("desa-ngebruk", uniqueFilename, file)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to AWS S3", "details": err.Error()})
			}

			// Membuat entitas gambar baru
			newImage := models.Image{
				OriginalName: file.Filename,
				Path:         s3Path,
				BlogID:       existingBlogger.Id_Blogger,
			}

			// Menyimpan gambar baru ke database
			if err := database.DB.Create(&newImage).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image to the database", "details": err.Error()})
			}
			existingBlogger.Images = append(existingBlogger.Images, newImage)
		}

	}

	existingBlogger.Name_Blog = updatedBlogger.Name_Blog
	existingBlogger.FillBlogger = updatedBlogger.FillBlogger
	existingBlogger.KategoriID = uint(kategoriIDInt)

	if err := database.DB.Model(&existingBlogger).Omit("CreatedAt").Save(&existingBlogger).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update blogger", "details": err.Error()})
	}

	response := fiber.Map{
		"message":  "Blogger updated successfully",
		"user_id":  userID,
		"username": username,
		"blogger":  existingBlogger,
	}

	return c.JSON(response)
}

// GetAllBloggers mengambil semua data blogger beserta informasi gambar dari database
func GetAllBloggers(c *fiber.Ctx) error {
	var bloggers []models.Blogger
	var total int64

	// Mengambil parameter page dan pageSize dari query string
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page size"})
	}

	// Mengambil parameter startDate dan endDate dari query string
	startDate := c.Query("startDate", "")
	endDate := c.Query("endDate", "")

	// Hitung offset berdasarkan page dan pageSize
	offset := (page - 1) * pageSize

	searchQuery := c.Query("search", "")
	selectedCategory := c.Query("category", "")

	// Membangun query dengan pencarian (jika ada)
	queryBuilder := database.DB.Model(&models.Blogger{}).
		Joins("left join kategoris on kategoris.id_kategori = bloggers.kategori_id").
		Joins("left join users on users.id_user = bloggers.user_id")
	if searchQuery != "" {
		queryBuilder = queryBuilder.Where("bloggers.name_blog LIKE ? OR kategoris.kategori_name LIKE ? OR users.nama LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Menambahkan filter berdasarkan tanggal (jika startDate dan endDate tidak kosong)
	if startDate != "" && endDate != "" {
		endDateWithTime := endDate + " 23:59:59"
		queryBuilder = queryBuilder.Where("bloggers.created_at BETWEEN ? AND ?", startDate, endDateWithTime)
	}

	if selectedCategory != "" {
		queryBuilder = queryBuilder.Where("kategoris.id_kategori = ?", selectedCategory)
	}

	// Menghitung jumlah total data (dengan pencarian dan filter tanggal, jika diterapkan)
	queryBuilder.Count(&total)

	// Mengambil data blogger dari database dengan pagination
	if err := queryBuilder.Preload("Images").Preload("Kategori").Preload("User").
		Offset(offset).Limit(pageSize).Find(&bloggers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bloggers", "details": err.Error()})
	}

	// Presign URL untuk setiap gambar di setiap blogger
	for i := range bloggers {
		// Handle heading
		if bloggers[i].Heading_Blogger != "" {
			headingPresignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", bloggers[i].Heading_Blogger)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL for heading", "details": err.Error()})
			}
			bloggers[i].Heading_Blogger = headingPresignURL
		}
	}
	for i := range bloggers {
		for j := range bloggers[i].Images {
			// Mendapatkan presign URL dari AWS S3
			presignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", bloggers[i].Images[j].Path)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL", "details": err.Error()})
			}

			// Mengganti path gambar dengan presign URL
			bloggers[i].Images[j].Path = presignURL
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := fiber.Map{
		"message":     "All bloggers retrieved successfully",
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalItems":  total,
		"bloggers":    bloggers,
	}

	return c.JSON(response)
}

// GetBloggerByID mengambil data blogger berdasarkan ID blogger
func GetBloggerByID(c *fiber.Ctx) error {
	id := c.Params("id_blogger")
	var blogger models.Blogger

	// Ambil data blogger berdasarkan ID
	if err := database.DB.Preload("Images").Preload("Kategori").Preload("User").First(&blogger, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blogger not found"})
	}

	// Generate presign URL untuk heading jika ada
	if blogger.Heading_Blogger != "" {
		headingPresignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", blogger.Heading_Blogger)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL for heading", "details": err.Error()})
		}
		blogger.Heading_Blogger = headingPresignURL
	}

	// Generate presign URL untuk setiap gambar
	for i, image := range blogger.Images {
		presignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", image.Path)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL", "details": err.Error()})
		}
		blogger.Images[i].Path = presignURL
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

	// Filter bloggers by user ID and retrieve data from database with related Images, Categories, and User
	if err := database.DB.Where("user_id = ?", userID).Preload("Images").Preload("Kategori").Preload("User").Find(&bloggers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bloggers for the specified user", "details": err.Error()})
	}

	// Presign URL for each image in each blogger
	for i := range bloggers {
		// Handle heading
		if bloggers[i].Heading_Blogger != "" {
			headingPresignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", bloggers[i].Heading_Blogger)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL for heading", "details": err.Error()})
			}
			bloggers[i].Heading_Blogger = headingPresignURL
		}

		for j := range bloggers[i].Images {
			// Generate presign URL from AWS S3
			presignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", bloggers[i].Images[j].Path)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL", "details": err.Error()})
			}

			// Replace image path with presign URL
			bloggers[i].Images[j].Path = presignURL
		}
	}

	response := fiber.Map{
		"message":  "Bloggers for the specified user retrieved successfully",
		"bloggers": bloggers,
	}

	return c.JSON(response)
}

func UserBloggersHandler(c *fiber.Ctx) error {
	// Extract the JWT token from the Authorization header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization token provided"})
	}

	// Verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Here, provide the key used to sign the tokens
		return []byte("Agnar123"), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Extract user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Handling different possible types for user ID
	var userID int
	if id, ok := claims["id_user"].(float64); ok { // If stored as a float64
		userID = int(id)
	} else if id, ok := claims["id_user"].(string); ok { // If stored as a string
		var err error
		userID, err = strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID in token is not an integer"})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
	}

	// Extract page and pageSize from query string
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page size"})
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Fetch bloggers for the specified user with pagination
	bloggers, total, err := getPaginatedBloggersByUserID(userID, offset, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve bloggers for the specified user", "details": err.Error()})
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Add the pre-signed heading URL to each blogger
	for i := range bloggers {
		if bloggers[i].Heading_Blogger != "" {
			headingPresignURL, err := awsutils.GeneratePresignURL("desa-ngebruk", bloggers[i].Heading_Blogger)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate presign URL for heading", "details": err.Error()})
			}
			bloggers[i].Heading_Blogger = headingPresignURL
		}
	}

	response := fiber.Map{
		"message":     "Bloggers for the specified user retrieved successfully",
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalItems":  total,
		"bloggers":    bloggers,
	}

	return c.JSON(response)
}

func getPaginatedBloggersByUserID(userID, offset, pageSize int) ([]models.Blogger, int, error) {
	var bloggers []models.Blogger
	var total int64

	// Filter bloggers by user ID and retrieve data from database with related Images, Categories, and User
	if err := database.DB.Where("user_id = ?", userID).Preload("Images").Preload("Kategori").Preload("User").
		Offset(offset).Limit(pageSize).Find(&bloggers).Error; err != nil {
		return nil, 0, err
	}

	// Count total number of bloggers for the specified user
	if err := database.DB.Model(&models.Blogger{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return bloggers, int(total), nil
}

func DeleteBlogger(c *fiber.Ctx) error {
	// Mendapatkan token JWT dari header Authorization
	tokenString := c.Get("Authorization")[7:]
	jwtToken := strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Mendapatkan ID Blogger dari parameter rute
	bloggerID := c.Params("id_blogger")

	// Mengambil informasi pengguna dari token JWT
	userID := uint(claims["id_user"].(float64))
	role := claims["role"].(string)

	var blogger models.Blogger
	if err := database.DB.Where("id_blogger = ?", bloggerID).First(&blogger).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blogger not found"})
	}

	// Validasi peran pengguna untuk mengizinkan atau menolak akses
	if !(role == "admin" || (role == "writter" && blogger.User_Id == userID)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	// Hapus blogger dari database
	if err := database.DB.Delete(&blogger).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete blogger from database", "details": err.Error()})
	}

	response := fiber.Map{
		"message": "Blogger deleted successfully",
		"user_id": userID,
		"blogger": blogger,
	}

	return c.JSON(response)
}

// Fungsi untuk menghapus gambar dari AWS S3
func deleteImageFromS3(bucket, objectKey string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	svc := s3.NewFromConfig(cfg)

	_, err = svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})

	return err
}

func uploadImageToS3(bucket, objectKey string, file *multipart.FileHeader) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	// Open the file
	srcFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Get file size
	fileSize, err := srcFile.Seek(0, io.SeekEnd)
	if err != nil {
		return "", err
	}

	// Reset file position for reading
	_, err = srcFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	// Prepare the upload input parameters
	uploadInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectKey),
		Body:          srcFile,
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(file.Header.Get("Content-Type")),
	}

	// Upload the file to S3
	_, err = client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		return "", err
	}

	// Mengembalikan hanya nama file unik sebagai path
	fmt.Printf("File uploaded successfully to S3. Path: %s\n", objectKey)

	return objectKey, nil
}

// CustomClaims holds custom claims for JWT.
type CustomClaims struct {
	UserID uint `json:"id_user"`
	jwt.StandardClaims
}

// extractUserClaimsFromToken extracts user claims from the JWT token.
func extractUserClaimsFromToken(c *fiber.Ctx) (*jwt.Token, *CustomClaims, error) {
	tokenString := c.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil // Replace with your JWT secret key
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

func getConfig() *aws.Config {
	cfg, err := awsutils.LoadAWSConfig()
	if err != nil {
		panic("configuration error: " + err.Error())
	}
	return cfg
}

func isValidImageFormat(file *multipart.FileHeader) bool {
	allowedFormats := map[string]bool{"image/png": true, "image/jpeg": true, "image/jpg": true}

	// Open the file to check its format
	f, err := file.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	// Determine the file format using the mimetype library
	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return false
	}

	// Convert the detected MIME type to lowercase for case-insensitive comparison
	detectedMIME := strings.ToLower(mime.String())

	return allowedFormats[detectedMIME]
}
func generateUniqueFilename(originalFilename string) string {
	// Generate a timestamp string to make the filename unique
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	// Extract the file extension
	fileExt := filepath.Ext(originalFilename)
	// Generate a unique filename by combining timestamp and original file extension
	return fmt.Sprintf("%d%s", timestamp, fileExt)
}
