package user

import (
	"gotickets/internal/auth"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService("")
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userHandler.CreateUser)
	api.POST("/login", userHandler.LoginUser)

}
