package event

import (
	"gotickets/internal/config"
	"log"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	// dependency injection
	// Initialize Cloudinary
	cloudinaryService, err := config.NewCloudinaryService()
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
	repo := NewRepository(db)
	service := NewService(repo, cloudinaryService)
	handler := NewHandler(service)

	api := e.Group("/api/v1/events")

	api.GET("", handler.GetEvents)
	api.GET("/:id", handler.GetEventById)
	api.POST("/create", handler.CreateEvent)
	api.PATCH("/:id", handler.UpdateEvent)
}
