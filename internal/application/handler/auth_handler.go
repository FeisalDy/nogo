package handler

import (
	"net/http"

	"github.com/FeisalDy/nogo/internal/application/service"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	userDto "github.com/FeisalDy/nogo/internal/user/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication requests at the application layer
type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// Register handles user registration with automatic default role assignment
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var registerDTO userDto.RegisterUserDTO
	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
		return
	}

	if err := h.validator.Struct(registerDTO); err != nil {
		utils.RespondValidationError(c, err, errors.ErrCodeUserValidation)
		return
	}

	if registerDTO.Password != registerDTO.ConfirmPassword {
		utils.RespondWithAppError(c, errors.ErrAuthPasswordMismatch)
		return
	}

	user, err := h.authService.Register(&registerDTO)
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

	response := userDto.AuthResponseDTO{
		Token: token,
		User: userDto.UserResponseDTO{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Bio:       user.Bio,
			Status:    user.Status,
		},
	}

	utils.RespondSuccess(c, http.StatusCreated, response, "Registration successful")
}
