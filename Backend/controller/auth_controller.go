package controller

import (
	"math"
	"strconv"
	"time"

	"github.com/agnarbriantama/DesaNgembruk-Backend/database"
	"github.com/agnarbriantama/DesaNgembruk-Backend/models"
	"github.com/agnarbriantama/DesaNgembruk-Backend/models/request"
	"github.com/agnarbriantama/DesaNgembruk-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validasi apakah email atau username sudah digunakan
	if isDuplicate := utils.IsDuplicateUser(user.Email, user.Username); isDuplicate {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email or username already exists"})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Set hashed password
	user.Password = string(hashedPassword)

	// Set default role jika tidak diisi
	if user.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Role is required"})
	}

	// Cek apakah nilai role yang diinginkan adalah "admin" atau "jobseeker"
	if user.Role != "admin" && user.Role != "writter" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role value"})
	}

	// Set waktu pembuatan dan pembaruan
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert user into the database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func Login(c *fiber.Ctx) error {
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Query user from the database by username
	var user models.User
	if err := database.DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Return user data and JWT token in the response
	responseData := fiber.Map{
		"user": fiber.Map{
			"id_user":    user.Id_User,
			"nama":       user.Nama,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			"role":       user.Role,
		},
		"token": token,
	}

	return c.JSON(responseData)
}

func UpdateUser(c *fiber.Ctx) error {
	// Validate JWT token
	// if err := utils.ProtectWithJWT(c); err != nil {
	//     return err
	// }

	id, err := strconv.Atoi(c.Params("id_user"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var updatedUserData models.User
	if err := c.BodyParser(&updatedUserData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Check if the username or email already exists in other user records
	var existingUser models.User
	database.DB.Where("username = ? AND id_user <> ?", updatedUserData.Username, id).
		Or("email = ? AND id_user <> ?", updatedUserData.Email, id).
		First(&existingUser)
	if existingUser.Id_User != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or email already in use by another user"})
	}

	// Encrypt password using bcrypt before updating
	if updatedUserData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUserData.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encrypt password"})
		}
		updatedUserData.Password = string(hashedPassword)
	}

	// Update user data
	result := database.DB.Model(&models.User{}).Where("id_user = ?", id).Updates(updatedUserData)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func UpdateUserWithoutPassword(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id_user"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var updatedUserData models.User
	if err := c.BodyParser(&updatedUserData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Check if username or email already exists
	var existingUser models.User
	database.DB.Where("username = ? AND id_user <> ?", updatedUserData.Username, id).Or("email = ? AND id_user <> ?", updatedUserData.Email, id).First(&existingUser)
	if existingUser.Id_User != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or email already in use"})
	}

	// Prevent password updates
	updatedUserData.Password = ""
	updatedUserData.Role = ""

	// Update user data without affecting the password and role
	result := database.DB.Model(&models.User{}).Where("id_user = ?", id).Omit("password", "role").Updates(updatedUserData)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// Fetch the updated user details
	var user models.User
	database.DB.First(&user, id)

	// Generate a new JWT token with the updated user data
	newToken, err := utils.GenerateJWTToken(user) // Note: Make sure this method accepts the updated user data
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new token"})
	}

	// Return the response with user details and new token
	responseData := fiber.Map{
		"user": fiber.Map{
			"id_user":    user.Id_User,
			"nama":       user.Nama,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			"role":       user.Role,
		},
		"token": newToken,
	}
	return c.JSON(responseData)
}

func GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "4"))

	var users []models.User
	var totalUsers int64

	database.DB.Model(&models.User{}).Count(&totalUsers)

	totalPages := int(math.Ceil(float64(totalUsers) / float64(pageSize)))

	offset := (page - 1) * pageSize
	result := database.DB.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	type paginatedUsersResponse struct {
		Users      []models.User `json:"users"`
		TotalPages int           `json:"totalPages"`
	}

	response := paginatedUsersResponse{
		Users:      users,
		TotalPages: totalPages,
	}

	return c.JSON(response)
}

func GetUserId(c *fiber.Ctx) error {
	idUser, err := strconv.Atoi(c.Params("id_user"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	result := database.DB.First(&user, "id_user = ?", idUser)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

func GetUserById(c *fiber.Ctx) error {

	userID, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	// Validate JWT token
	if err := utils.ProtectWithJWT(c, "admin"); err != nil {
		return err
	}
	id := c.Params("id_user")

	var user models.User
	if err := database.DB.Where("id_user = ?", id).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User successfully deleted"})
}

func GetBloggersByUser(c *fiber.Ctx) error {
	idUser := c.Params("id_user")

	var bloggers []models.Blogger
	result := database.DB.Where("user_id = ?", idUser).Find(&bloggers)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot retrieve bloggers"})
	}

	return c.JSON(bloggers)
}

func ChangePassword(c *fiber.Ctx) error {
	id_user, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	var req request.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	result := database.DB.First(&user, id_user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Verifikasi old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Old password is incorrect"})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash password"})
	}

	// Update password di database
	user.Password = string(hashedPassword)
	database.DB.Save(&user)

	return c.JSON(fiber.Map{"message": "Password successfully changed"})
}
