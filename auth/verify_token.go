package auth

import (
	"fmt"
	"net/http"
	"os"
	"todolist/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func ExtractUsername(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
		}

		tokenString := ""
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(*models.JWTClaims)
		if !ok || claims == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid claims"})
		}

		c.Set("username", claims.Username)

		return next(c)
	}
}
