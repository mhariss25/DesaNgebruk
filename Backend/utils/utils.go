package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"DesaNgebruk/database"
	"DesaNgebruk/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func IsDuplicateUser(email, username string) bool {
	var existingUser models.User
	// Cek apakah email sudah ada
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return true
	}
	// Cek apakah username sudah ada
	if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return true
	}
	return false
}

// SecretKey digunakan untuk menandatangani token JWT
var SecretKey = []byte("Agnar123")

// GenerateJWTToken membuat token JWT berdasarkan informasi pengguna
func GenerateJWTToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id_user"] = user.Id_User
	claims["UserID"] = user.Id_User
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// middleware/middleware.go

// func ProtectWithJWT(c *fiber.Ctx) error {
// 	// Mendapatkan token dari header Authorization
// 	tokenString := c.Get("Authorization")

// 	// Memeriksa keberadaan token
// 	if tokenString == "" {
// 		return fiber.ErrUnauthorized
// 	}

// 	// Validasi token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		jwtSecret := os.Getenv("JWT_SECRET")
// 		return []byte(jwtSecret), nil // Ganti dengan secret key yang sesuai
// 	})

// 	if err != nil || !token.Valid {
// 		return fiber.ErrUnauthorized
// 	}

// 	// Mendapatkan role dari token
// 	userClaims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok || userClaims == nil {
// 		return fiber.ErrUnauthorized
// 	}

// 	userRole, ok := userClaims["role"].(string)
// 	if !ok || userRole == "" {
// 		return fiber.ErrUnauthorized
// 	}

// 	// Memeriksa apakah role pengguna termasuk dalam role yang diizinkan
// 	roleAllowed := false
// 	for _, allowedRole := range []string{"jobseeker"} {
// 		if userRole == allowedRole {
// 			roleAllowed = true
// 			break
// 		}
// 	}

// 	if !roleAllowed {
// 		return fiber.ErrUnauthorized
// 	}

// 	// Menambahkan informasi token ke konteks dengan kunci "user"
// 	c.Locals("user", userClaims)

// 	return nil
// }

func ProtectWithJWT(c *fiber.Ctx, allowedRoles ...string) error {
	// Mendapatkan token dari header Authorization
	tokenString := c.Get("Authorization")

	// Memeriksa keberadaan token
	if tokenString == "" {
		return fiber.ErrUnauthorized
	}

	// Validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("Agnar123"), nil // Ganti dengan secret key yang sesuai
	})

	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	// Mendapatkan role dari token
	userClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || userClaims == nil {
		return fiber.ErrUnauthorized
	}

	userRole, ok := userClaims["role"].(string)
	if !ok || userRole == "" {
		return fiber.ErrUnauthorized
	}

	// Memeriksa apakah role pengguna termasuk dalam role yang diizinkan
	roleAllowed := false
	for _, allowedRole := range allowedRoles {
		if userRole == allowedRole {
			roleAllowed = true
			break
		}
	}

	if !roleAllowed {
		return fiber.ErrUnauthorized
	}

	// Menambahkan informasi token ke konteks dengan kunci "user"
	c.Locals("user", userClaims)

	return nil
}

// Fungsi untuk mendapatkan id_user dari token JWT
func GetUserIdFromToken(c *fiber.Ctx) (uint, error) {
	token := c.Get("Authorization")

	// Mendapatkan token tanpa "Bearer "
	jwtToken := strings.TrimPrefix(token, "Bearer ")

	// Parse token JWT
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if !parsedToken.Valid {
		return 0, errors.New("Invalid token")
	}

	// Mendapatkan id_user dari claim JWT
	idUserFloat64, ok := claims["id_user"].(float64)
	if !ok {
		return 0, errors.New("Invalid id_user claim in token")
	}

	idUser := uint(idUserFloat64)
	return idUser, nil
}
