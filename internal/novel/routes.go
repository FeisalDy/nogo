package novel

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all novel-related routes
func RegisterRoutes(router *gin.RouterGroup) {
	// Placeholder routes for novel domain
	// You can implement these as you build the novel functionality

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Novel API endpoint",
			"status":  "coming soon",
		})
	})

	// Future novel routes:
	// router.POST("/", novelHandler.CreateNovel)
	// router.GET("/:id", novelHandler.GetNovel)
	// router.PUT("/:id", novelHandler.UpdateNovel)
	// router.DELETE("/:id", novelHandler.DeleteNovel)
	// router.GET("/", novelHandler.GetAllNovels)
	// router.GET("/:id/chapters", novelHandler.GetNovelChapters)
}
