package novel

import (
	"github.com/FeisalDy/nogo/internal/novel/handler"
	"github.com/FeisalDy/nogo/internal/novel/repository"
	"github.com/FeisalDy/nogo/internal/novel/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.RouterGroup) {
	novelRepo := repository.NewNovelRepository(db)
	novelSrvc := service.NewNovelService(novelRepo)
	novelHandler := handler.NewNovelHandler(novelSrvc)

	novelRoutes := router.Group("/")
	{
		// Single novel operations
		novelRoutes.GET("/:id", novelHandler.GetNovelByID)

		// Cursor-based pagination endpoints
		novelRoutes.GET("", novelHandler.GetAllNovels) // GET /novels?cursor=...&limit=20
	}
}
