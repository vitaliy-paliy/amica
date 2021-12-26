package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vitaliy-paliy/amica/pkg/models"
)

type jwtAuthClaims struct {
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

type jwtUserClaims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

func SetAuthMiddleware(g *echo.Group) {
	config := middleware.JWTConfig{
		Claims:     &jwtAuthClaims{},
		SigningKey: []byte("auth-jwt-secret"),
	}

	g.Use(middleware.JWTWithConfig(config))
}

func SetUserMiddleware(g *echo.Group) {
	config := middleware.JWTConfig{
		Claims:     &jwtUserClaims{},
		SigningKey: []byte("user-jwt-secret"),
	}

	g.Use(middleware.JWTWithConfig(config))
}
