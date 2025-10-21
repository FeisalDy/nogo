package handler

import (
	"net/http"

	"github.com/FeisalDy/nogo/internal/application/service"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/gin-gonic/gin"
)

// UserProfileHandler handles user profile requests at the application layer
// This is for operations that span multiple domains (User + Role + Casbin)
type UserProfileHandler struct {
	userProfileService *service.UserProfileService
}

// NewUserProfileHandler creates a new instance of UserProfileHandler
func NewUserProfileHandler(userProfileService *service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileService: userProfileService,
	}
}

// GetMe retrieves the current authenticated user's profile with roles and permissions
// GET /api/v1/profile/me
func (h *UserProfileHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
		return
	}

	id := userID.(uint)
	userWithPermissions, err := h.userProfileService.GetUserWithPermissions(id)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, userWithPermissions, "User profile retrieved successfully")
}
