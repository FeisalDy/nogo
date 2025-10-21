package dto

type RegisterUserDTO struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponseDTO struct {
	ID        uint    `json:"id"`
	Username  *string `json:"username"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Status    string  `json:"status"`
}

type RoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PermissionDTO struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type UserWithPermissionsDTO struct {
	ID          uint            `json:"id"`
	Username    *string         `json:"username"`
	Email       string          `json:"email"`
	AvatarURL   *string         `json:"avatar_url,omitempty"`
	Bio         *string         `json:"bio,omitempty"`
	Status      string          `json:"status"`
	Roles       []RoleDTO       `json:"roles"`
	Permissions []PermissionDTO `json:"permissions"`
}

type AuthResponseDTO struct {
	Token string          `json:"token"`
	User  UserResponseDTO `json:"user"`
}
