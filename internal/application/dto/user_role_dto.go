package dto

// AssignRoleToUserDTO represents the request to assign a role to a user
type AssignRoleToUserDTO struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

// RemoveRoleFromUserDTO represents the request to remove a role from a user
type RemoveRoleFromUserDTO struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

// UserRoleAssignmentDTO represents a user-role assignment
type UserRoleAssignmentDTO struct {
	UserID    uint   `json:"user_id"`
	RoleID    uint   `json:"role_id"`
	RoleName  string `json:"role_name"`
	UserEmail string `json:"user_email,omitempty"`
	Username  string `json:"username,omitempty"`
}
