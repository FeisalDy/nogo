package common

import (
	"boiler/internal/common/handler"
	"boiler/internal/common/service"

	"github.com/gin-gonic/gin"
)

// RegisterUploadRoutes registers all upload-related routes
func RegisterUploadRoutes(router *gin.RouterGroup, baseURL string) {
	// Initialize upload service
	uploadConfig := service.UploadConfig{
		UploadDir:    "./uploads", // Directory to store files
		BaseURL:      baseURL,     // Use configured base URL
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
		MaxSize:      10 * 1024 * 1024, // 10MB
	}
	uploadService := service.NewUploadService(uploadConfig)

	// Upload handler config
	uploadHandlerConfig := handler.UploadHandlerConfig{
		DefaultAllowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		DefaultMaxSize:      5 * 1024 * 1024, // 5MB default
	}
	uploadHandler := handler.NewUploadHandler(uploadService, uploadHandlerConfig)

	// Register upload routes
	router.POST("/single", uploadHandler.UploadSingleFile)
	router.POST("/multiple", uploadHandler.UploadMultipleFiles)
	router.DELETE("/file", uploadHandler.DeleteFile)

	// Add more upload routes as needed:
	// router.GET("/files", uploadHandler.ListFiles)
	// router.GET("/file/:filename/info", uploadHandler.GetFileInfo)
}

// GetUploadService returns a configured upload service for use by other domains
func GetUploadService(baseURL string) *service.UploadService {
	uploadConfig := service.UploadConfig{
		UploadDir:    "./uploads",
		BaseURL:      baseURL,
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
		MaxSize:      10 * 1024 * 1024, // 10MB
	}
	return service.NewUploadService(uploadConfig)
}
