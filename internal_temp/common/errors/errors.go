package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// AppError represents a custom application error with code and message
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// Error codes for different domains
const (
	// User domain errors (USER001-USER099)
	ErrCodeUserNotFound           = "USER001"
	ErrCodeUserAlreadyExists      = "USER002"
	ErrCodeUserCreationFailed     = "USER003"
	ErrCodeUserUpdateFailed       = "USER004"
	ErrCodeUserDeletionFailed     = "USER005"
	ErrCodeUserInvalidCredentials = "USER006"
	ErrCodeUserValidation         = "USER007"

	// Role domain errors (ROLE001-ROLE099)
	ErrCodeRoleNotFound       = "ROLE001"
	ErrCodeRoleAlreadyExists  = "ROLE002"
	ErrCodeRoleCreationFailed = "ROLE003"
	ErrCodeRoleUpdateFailed   = "ROLE004"
	ErrCodeRoleDeletionFailed = "ROLE005"
	ErrCodeRoleValidation     = "ROLE006"

	// User-Role domain errors (USERROLE001-USERROLE099)
	ErrCodeUserRoleNotFound       = "USERROLE001"
	ErrCodeUserRoleAlreadyExists  = "USERROLE002"
	ErrCodeUserRoleCreationFailed = "USERROLE003"
	ErrCodeUserRoleUpdateFailed   = "USERROLE004"
	ErrCodeUserRoleDeletionFailed = "USERROLE005"

	// Auth domain errors (AUTH001-AUTH099)
	ErrCodeAuthInvalidToken       = "AUTH001"
	ErrCodeAuthTokenExpired       = "AUTH002"
	ErrCodeAuthTokenMissing       = "AUTH003"
	ErrCodeAuthUnauthorized       = "AUTH004"
	ErrCodeAuthForbidden          = "AUTH005"
	ErrCodeAuthPasswordMismatch   = "AUTH006"
	ErrCodeAuthRegistrationFailed = "AUTH007"
	ErrCodeAuthLoginFailed        = "AUTH008"

	// Upload domain errors (UPLOAD001-UPLOAD099)
	ErrCodeUploadInvalidFile  = "UPLOAD001"
	ErrCodeUploadFailed       = "UPLOAD002"
	ErrCodeUploadFileTooLarge = "UPLOAD003"
	ErrCodeUploadInvalidType  = "UPLOAD004"
	ErrCodeUploadNoFile       = "UPLOAD005"

	// Database errors (DB001-DB099)
	ErrCodeDatabaseConnection  = "DB001"
	ErrCodeDatabaseQuery       = "DB002"
	ErrCodeDatabaseTransaction = "DB003"

	// Validation errors (VAL001-VAL099)
	ErrCodeValidationFailed = "VAL001"
	ErrCodeInvalidInput     = "VAL002"
	ErrCodeMissingField     = "VAL003"
	ErrCodeInvalidParam     = "VAL004"

	// Casbin errors (CASBIN001-CASBIN099)
	ErrCodeCasbinPolicyLoadFailed   = "CASBIN001"
	ErrCodeCasbinPolicySaveFailed   = "CASBIN002"
	ErrCodeCasbinPolicyRemoveFailed = "CASBIN003"

	// General errors (GEN001-GEN099)
	ErrCodeInternalServer = "GEN001"
	ErrCodeBadRequest     = "GEN002"
	ErrCodeNotFound       = "GEN003"

	// File Error
)

var (
	// user related
	ErrUserNotFound           = NewAppError(ErrCodeUserNotFound, "User not found")
	ErrUserAlreadyExists      = NewAppError(ErrCodeUserAlreadyExists, "User already exists")
	ErrUserCreationFailed     = NewAppError(ErrCodeUserCreationFailed, "Failed to create user")
	ErrUserUpdateFailed       = NewAppError(ErrCodeUserUpdateFailed, "Failed to update user")
	ErrUserDeletionFailed     = NewAppError(ErrCodeUserDeletionFailed, "Failed to delete user")
	ErrUserInvalidCredentials = NewAppError(ErrCodeUserInvalidCredentials, "Invalid username or password")

	//user - role related
	ErrUserRoleAssignmentFailed = NewAppError(ErrCodeUserRoleNotFound, "Failed to assign role to user")
	ErrUserRoleRemovalFailed    = NewAppError(ErrCodeUserRoleAlreadyExists, "Failed to remove role from user")
	ErrUserDoesNotHaveRole      = NewAppError(ErrCodeUserRoleCreationFailed, "User does not have the specified role")
	ErrUserAlreadyHasRole       = NewAppError(ErrCodeUserRoleUpdateFailed, "User already has the specified role")

	// role related
	ErrRoleNotFound       = NewAppError(ErrCodeRoleNotFound, "Role not found")
	ErrRoleAlreadyExists  = NewAppError(ErrCodeRoleAlreadyExists, "Role already exists")
	ErrRoleCreationFailed = NewAppError(ErrCodeRoleCreationFailed, "Failed to create role")
	ErrRoleUpdateFailed   = NewAppError(ErrCodeRoleUpdateFailed, "Failed to update role")
	ErrRoleDeletionFailed = NewAppError(ErrCodeRoleDeletionFailed, "Failed to delete role")

	// auth related
	ErrAuthInvalidToken     = NewAppError(ErrCodeAuthInvalidToken, "Invalid authentication token")
	ErrAuthTokenExpired     = NewAppError(ErrCodeAuthTokenExpired, "Authentication token has expired")
	ErrAuthTokenMissing     = NewAppError(ErrCodeAuthTokenMissing, "Authentication token is missing")
	ErrAuthUnauthorized     = NewAppError(ErrCodeAuthUnauthorized, "Unauthorized access")
	ErrAuthForbidden        = NewAppError(ErrCodeAuthForbidden, "Access forbidden")
	ErrAuthPasswordMismatch = NewAppError(ErrCodeAuthPasswordMismatch, "Password and confirm password do not match")
	ErrAuthLoginFailed      = NewAppError(ErrCodeAuthLoginFailed, "Login failed")

	// upload related
	ErrUploadInvalidFile  = NewAppError(ErrCodeUploadInvalidFile, "Invalid file")
	ErrUploadFailed       = NewAppError(ErrCodeUploadFailed, "Failed to upload file")
	ErrUploadFileTooLarge = NewAppError(ErrCodeUploadFileTooLarge, "File size exceeds maximum limit")
	ErrUploadInvalidType  = NewAppError(ErrCodeUploadInvalidType, "Invalid file type")
	ErrUploadNoFile       = NewAppError(ErrCodeUploadNoFile, "No file provided")

	// database related
	ErrDatabaseConnection  = NewAppError(ErrCodeDatabaseConnection, "Database connection error")
	ErrDatabaseQuery       = NewAppError(ErrCodeDatabaseQuery, "Database query error")
	ErrDatabaseTransaction = NewAppError(ErrCodeDatabaseTransaction, "Database transaction error")

	// validation related
	ErrValidationFailed = NewAppError(ErrCodeValidationFailed, "Validation failed")
	ErrInvalidInput     = NewAppError(ErrCodeInvalidInput, "Invalid input")
	ErrMissingField     = NewAppError(ErrCodeMissingField, "Required field is missing")
	ErrInvalidParam     = NewAppError(ErrCodeInvalidParam, "Invalid parameter")

	// casbin related
	ErrCasbinPolicyLoadFailed   = NewAppError(ErrCodeCasbinPolicyLoadFailed, "Failed to load Casbin policies")
	ErrCasbinPolicySaveFailed   = NewAppError(ErrCodeCasbinPolicySaveFailed, "Failed to save Casbin policies")
	ErrCasbinPolicyRemoveFailed = NewAppError(ErrCodeCasbinPolicyRemoveFailed, "Failed to remove Casbin policies")

	// general related
	ErrInternalServer = NewAppError(ErrCodeInternalServer, "Internal server error")
	ErrBadRequest     = NewAppError(ErrCodeBadRequest, "Bad request")
	ErrNotFound       = NewAppError(ErrCodeNotFound, "Resource not found")
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

// FormatValidationError formats validator errors into readable messages
func FormatValidationError(err error, errorCode string) *AppError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return NewAppError(errorCode, err.Error())
	}

	var errors []ValidationError
	for _, fieldError := range validationErrors {
		field := fieldError.Field()
		tag := fieldError.Tag()

		// Convert field name from PascalCase to snake_case or readable format
		readableField := toReadableField(field)

		message := getValidationMessage(fieldError)

		errors = append(errors, ValidationError{
			Field:   readableField,
			Message: message,
			Tag:     tag,
			Value:   fmt.Sprintf("%v", fieldError.Value()),
		})
	}

	// Create a human-readable summary message
	var messages []string
	for _, e := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}

	appError := NewAppError(errorCode, "Validation failed")
	appError.Details = map[string]any{
		"fields":  errors,
		"summary": strings.Join(messages, "; "),
	}

	return appError
}

// getValidationMessage returns a human-readable message for a validation error
func getValidationMessage(fe validator.FieldError) string {
	field := toReadableField(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, toReadableField(fe.Param()))
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be a valid number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s failed validation on %s", field, fe.Tag())
	}
}

// toReadableField converts field names to a more readable format
// Example: "ConfirmPassword" -> "confirm password"
func toReadableField(field string) string {
	if field == "" {
		return field
	}

	var result []rune
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, ' ')
		}
		if i == 0 {
			result = append(result, r)
		} else {
			result = append(result, r)
		}
	}

	return strings.ToLower(string(result))
}
