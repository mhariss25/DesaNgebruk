package controller

import (
	"context"
	"time"

	"github.com/agnarbriantama/DesaNgembruk-Backend/database"
	"github.com/agnarbriantama/DesaNgembruk-Backend/models"
	"github.com/agnarbriantama/DesaNgembruk-Backend/utils/awsutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

func UploadImages(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-2")) // Ganti dengan region AWS Anda
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(cfg)

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["images"]
	presignedURLs := make([]string, 0)

	for _, file := range files {
		// Proses setiap file
		content, err := file.Open()
		if err != nil {
			return err
		}
		defer content.Close()

		uniqueFilename := generateUniqueFilename(file.Filename)
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("desa-ngebruk"),
			Key:    aws.String(uniqueFilename),
			Body:   content,
		})
		if err != nil {
			return err
		}

		// Generate presigned URL
		presignedURL, err := awsutils.GeneratePublicURL("desa-ngebruk", uniqueFilename)
		if err != nil {
			return err
		}

		// Simpan info gambar di MySQL
		img := models.Gambar{
			OriginalName: file.Filename,
			Path:         uniqueFilename,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := database.DB.Create(&img).Error; err != nil {
			return err
		}

		// Tambahkan presigned URL ke slice
		presignedURLs = append(presignedURLs, presignedURL)
	}

	// Berikan respons dengan presigned URLs
	return c.JSON(fiber.Map{"message": "Images uploaded successfully", "presigned_urls": presignedURLs})
}
