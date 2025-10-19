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

func (h *UserHandler) Register(c *gin.Context) {
	var registerDTO dto.RegisterUserDTO
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

	user, err := h.userService.Register(&registerDTO)
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

	utils.RespondSuccess(c, http.StatusCreated, response, "Registration successful")
}

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

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondWithAppError(c, errors.ErrAuthUnauthorized)
		return
	}

	id := userID.(uint)
	userWithPermissions, err := h.userService.GetUserWithPermissions(id)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondSuccess(c, http.StatusOK, userWithPermissions, "User profile retrieved successfully")
}

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
