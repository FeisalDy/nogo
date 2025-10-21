package dto

// RoleDTO represents role data for responses
type RoleDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// CreateRoleDTO for creating a new role
type CreateRoleDTO struct {
	Name        string  `json:"name" validate:"required,min=3,max=50"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateRoleDTO for updating a role
type UpdateRoleDTO struct {
	Name        *string `json:"name" validate:"omitempty,min=3,max=50"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// RoleWithUsersDTO represents a role with its users
type RoleWithUsersDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Users       []UserDTO `json:"users"`
	UserCount   int64     `json:"user_count"`
}

// UserDTO represents a user (from user domain)
type UserDTO struct {
	ID        uint    `json:"id"`
	Username  *string `json:"username"`
	Email     string  `json:"email"`
	Status    string  `json:"status"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
