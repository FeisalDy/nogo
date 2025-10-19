package utils

import (
	"net/http"

	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// RespondSuccess sends a success response
func RespondSuccess(c *gin.Context, statusCode int, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Data:    data,
		Message: msg,
	})
}

// RespondError sends an error response
func RespondError(c *gin.Context, statusCode int, appError *errors.AppError) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Code:    appError.Code,
		Message: appError.Message,
		Details: appError.Details,
	})
}

func HandleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if appErr, ok := err.(*errors.AppError); ok {
		RespondWithAppError(c, appErr)
	} else {
		appError := errors.ErrInternalServer.WithDetails(map[string]any{
			"reason": err.Error(),
		})
		RespondWithAppError(c, appError)
	}
}

// RespondValidationError sends a validation error response
func RespondValidationError(c *gin.Context, err error, errorCode string) {
	appError := errors.FormatValidationError(err, errorCode)
	RespondError(c, http.StatusBadRequest, appError)
}

// RespondWithAppError sends an error response with appropriate status code
func RespondWithAppError(c *gin.Context, appError *errors.AppError) {
	statusCode := GetStatusCodeFromErrorCode(appError.Code)
	RespondError(c, statusCode, appError)
}

// GetStatusCodeFromErrorCode maps error codes to HTTP status codes
func GetStatusCodeFromErrorCode(code string) int {
	switch code {
	// User errors
	case errors.ErrCodeUserNotFound:
		return http.StatusNotFound
	case errors.ErrCodeUserAlreadyExists:
		return http.StatusConflict
	case errors.ErrCodeUserCreationFailed, errors.ErrCodeUserUpdateFailed, errors.ErrCodeUserDeletionFailed:
		return http.StatusInternalServerError
	case errors.ErrCodeUserInvalidCredentials:
		return http.StatusUnauthorized
	case errors.ErrCodeUserValidation:
		return http.StatusBadRequest

	// Auth errors
	case errors.ErrCodeAuthInvalidToken, errors.ErrCodeAuthTokenExpired, errors.ErrCodeAuthTokenMissing, errors.ErrCodeAuthUnauthorized, errors.ErrCodeAuthLoginFailed:
		return http.StatusUnauthorized
	case errors.ErrCodeAuthForbidden:
		return http.StatusForbidden
	case errors.ErrCodeAuthPasswordMismatch, errors.ErrCodeAuthRegistrationFailed:
		return http.StatusBadRequest

	// Upload errors
	case errors.ErrCodeUploadInvalidFile, errors.ErrCodeUploadFileTooLarge, errors.ErrCodeUploadInvalidType, errors.ErrCodeUploadNoFile:
		return http.StatusBadRequest
	case errors.ErrCodeUploadFailed:
		return http.StatusInternalServerError

	// Database errors
	case errors.ErrCodeDatabaseConnection, errors.ErrCodeDatabaseQuery, errors.ErrCodeDatabaseTransaction:
		return http.StatusInternalServerError

	// Validation errors
	case errors.ErrCodeValidationFailed, errors.ErrCodeInvalidInput, errors.ErrCodeMissingField:
		return http.StatusBadRequest

	// General errors
	case errors.ErrCodeInternalServer:
		return http.StatusInternalServerError
	case errors.ErrCodeBadRequest:
		return http.StatusBadRequest
	case errors.ErrCodeNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
