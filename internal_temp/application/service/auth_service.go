package service

import (
	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/database"
	roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
	userDto "github.com/FeisalDy/nogo/internal/user/dto"
	userModel "github.com/FeisalDy/nogo/internal/user/model"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
	"gorm.io/gorm"
)

// AuthService handles authentication operations that span multiple domains
// This is part of the Application Layer
type AuthService struct {
	userRepo      *userRepo.UserRepository
	roleRepo      *roleRepo.RoleRepository
	casbinService *casbinService.CasbinService
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepository *userRepo.UserRepository,
	roleRepository *roleRepo.RoleRepository,
	casbin *casbinService.CasbinService,
) *AuthService {
	return &AuthService{
		userRepo:      userRepository,
		roleRepo:      roleRepository,
		casbinService: casbin,
	}
}

// Register handles user registration with default role assignment
// This is a cross-domain operation that:
// 1. Creates a new user (User domain)
// 2. Assigns default "user" role (Role domain)
// 3. Syncs with Casbin for authorization
func (s *AuthService) Register(registerDTO *userDto.RegisterUserDTO) (*userModel.User, error) {
	var userCreated *userModel.User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Check if user already exists
		existingUser, _ := s.userRepo.WithTx(tx).GetUserByEmail(registerDTO.Email)
		if existingUser != nil {
			return errors.ErrUserAlreadyExists
		}

		// 2. Hash password
		hashedPassword, err := utils.HashPassword(registerDTO.Password)
		if err != nil {
			return err
		}

		// 3. Create user
		user := &userModel.User{
			Username: &registerDTO.Username,
			Email:    registerDTO.Email,
			Password: &hashedPassword,
		}

		userCreated, err = s.userRepo.WithTx(tx).CreateUser(user)
		if err != nil {
			return err
		}

		// 4. Get default "user" role
		defaultRole, err := s.roleRepo.WithTx(tx).GetByName("user")
		if err != nil {
			return errors.ErrRoleNotFound
		}

		// 5. Assign role in database (user_roles table)
		if err := s.userRepo.WithTx(tx).AssignRoleToUser(userCreated.ID, defaultRole.ID); err != nil {
			return err
		}

		// 6. Assign role in Casbin (for authorization checks)
		// Note: This is done within the transaction context
		if err := s.casbinService.AssignRoleToUser(userCreated.ID, defaultRole.Name); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return userCreated, nil
}
