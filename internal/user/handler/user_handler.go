package handler

import (
	"net/http"

	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/user/dto"
	"github.com/FeisalDy/nogo/internal/user/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService, validator: validator.New()}
}

// Note: Register has been moved to the application layer (internal/application/handler/auth_handler.go)
// This is because registration involves cross-domain operations (creating user + assigning role)

func (h *UserHandler) Login(c *gin.Context) {
	var loginDTO dto.LoginUserDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
		return
	}

	if err := h.validator.Struct(loginDTO); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
		return
	}

	user, err := h.userService.Login(&loginDTO)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	username := ""
	if user.Username != nil {
		username = *user.Username
	}
	token, err := utils.GenerateToken(user.ID, user.Email, username)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Prepare response
	response := dto.AuthResponseDTO{
		Token: token,
		User: dto.UserResponseDTO{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Bio:       user.Bio,
			Status:    user.Status,
		},
	}

	utils.RespondSuccess(c, http.StatusOK, response, "Login successful")
}

// Note: GetMe has been moved to the application layer (internal/application/handler/user_profile_handler.go)
// This is because it involves cross-domain operations (user + roles + permissions from Casbin)
// New endpoint: GET /api/v1/profile/me

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	emailParam := c.Param("email")

	user, err := h.userService.GetUserByEmail(emailParam)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	res := dto.UserResponseDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Status:    user.Status,
	}

	utils.RespondSuccess(c, http.StatusOK, res)
}
