package user_handler

import (
	"net/http"
	"time"
	"todolist/auth"
	"todolist/helper"
	"todolist/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var DB *gorm.DB

func RegisterUserHandlers(db *gorm.DB) {
	DB = db
}

func Register(c echo.Context) error {
	var newUser models.User
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	var existingUser models.User
	if err := DB.Where("username = ? OR email = ?", newUser.Username, newUser.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username or Email already exists"})
	}

	newUser.ID = helper.GenerateRandomString(10)
	if err := DB.Create(&newUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func Login(c echo.Context) error {
	var credentials models.User
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	var user models.User
	if err := DB.Where("username = ? AND password = ?", credentials.Username, credentials.Password).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.JWTClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	tokenString, err := token.SignedString(auth.JwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}
