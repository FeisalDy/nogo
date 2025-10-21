package handler

import (
	"net/http"
	"strconv"

	"github.com/FeisalDy/nogo/internal/application/dto"
	"github.com/FeisalDy/nogo/internal/application/service"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserRoleHandler handles HTTP requests for user-role operations
// This is part of the Application Layer that handles cross-domain operations
type UserRoleHandler struct {
	userRoleService *service.UserRoleService
	validator       *validator.Validate
}

// NewUserRoleHandler creates a new instance of UserRoleHandler
func NewUserRoleHandler(userRoleService *service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{
		userRoleService: userRoleService,
		validator:       validator.New(),
	}
}

// AssignRoleToUser handles the request to assign a role to a user
// POST /api/v1/user-roles/assign
func (h *UserRoleHandler) AssignRoleToUser(c *gin.Context) {
	var req dto.AssignRoleToUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
		return
	}

	if err := h.userRoleService.AssignRoleToUser(req.UserID, req.RoleID); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"user_id": req.UserID,
		"role_id": req.RoleID,
	}, "Role assigned to user successfully")
}

// RemoveRoleFromUser handles the request to remove a role from a user
// POST /api/v1/user-roles/remove
func (h *UserRoleHandler) RemoveRoleFromUser(c *gin.Context) {
	var req dto.RemoveRoleFromUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeValidationFailed)
		return
	}

	if err := h.userRoleService.RemoveRoleFromUser(req.UserID, req.RoleID); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"user_id": req.UserID,
		"role_id": req.RoleID,
	}, "Role removed from user successfully")
}

// GetUserRoles handles the request to get all roles for a user
// GET /api/v1/user-roles/users/:user_id/roles
func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.RespondWithAppError(c, errors.ErrInvalidParam.WithDetails(map[string]any{
			"reason": err.Error(),
		}))
		return
	}

	roles, err := h.userRoleService.GetUserRoles(uint(userID))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, roles, "User roles retrieved successfully")
}
