package handler

import (
	"net/http"
	"strconv"

	commonDto "github.com/FeisalDy/nogo/internal/common/dto"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/novel/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type NovelHandler struct {
	novelService *service.NovelService
	validator    *validator.Validate
}

func NewNovelHandler(novelService *service.NovelService) *NovelHandler {
	return &NovelHandler{novelService: novelService, validator: validator.New()}
}

func (h *NovelHandler) GetNovelByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)

	if err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeInvalidParam)
		return
	}

	novel, err := h.novelService.GetNovelByID(uint(id))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, novel)
}

// GetAllNovels retrieves all novels with cursor-based pagination
// Query params:
//   - cursor: base64-encoded cursor for pagination (optional)
//   - limit: number of items per page (default: 20, max: 100)
//   - sort_order: "asc" or "desc" (default: "desc")
func (h *NovelHandler) GetAllNovels(c *gin.Context) {
	var req commonDto.CursorPaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeInvalidParam)
		return
	}

	novels, pageInfo, err := h.novelService.GetAllNovelsWithCursor(&req)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Use the new pagination response helper (no nested data)
	utils.RespondSuccessWithPagination(
		c,
		http.StatusOK,
		novels, // Direct data array
		pageInfo,
		commonDto.PaginationMetadata{
			Count:     len(novels),
			Limit:     req.Limit,
			SortOrder: req.SortOrder,
		},
	)
}
