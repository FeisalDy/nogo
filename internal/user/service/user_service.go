package service

import (
	"strings"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/common/utils"
	"github.com/FeisalDy/nogo/internal/database"
	roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
	"github.com/FeisalDy/nogo/internal/user/dto"
	"github.com/FeisalDy/nogo/internal/user/model"
	"github.com/FeisalDy/nogo/internal/user/repository"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo      *repository.UserRepository
	roleRepo      *roleRepo.RoleRepository
	casbinService *casbinService.CasbinService
}

func NewUserService(userRepository *repository.UserRepository, roleRepository *roleRepo.RoleRepository, casbin *casbinService.CasbinService) *UserService {
	return &UserService{
		userRepo:      userRepository,
		roleRepo:      roleRepository,
		casbinService: casbin,
	}
}

func (s *UserService) Register(registerDTO *dto.RegisterUserDTO) (*model.User, error) {
	var userCreated *model.User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		existingUser, _ := s.userRepo.WithTx(tx).GetUserByEmail(registerDTO.Email)
		if existingUser != nil {
			return errors.ErrUserAlreadyExists
		}

		hashedPassword, err := utils.HashPassword(registerDTO.Password)
		if err != nil {
			return err
		}

		user := &model.User{
			Username: &registerDTO.Username,
			Email:    registerDTO.Email,
			Password: &hashedPassword,
		}

		userCreated, err = s.userRepo.WithTx(tx).CreateUser(user)
		if err != nil {
			return err
		}

		defaultRole, err := s.roleRepo.WithTx(tx).GetByName("user")
		if err != nil {
			return errors.ErrRoleNotFound
		}

		// Assign role in database (user_roles table)
		if err := s.userRepo.WithTx(tx).AssignRoleToUser(userCreated.ID, defaultRole.ID); err != nil {
			return err
		}

		// Assign role in Casbin (for authorization checks)
		// Note: This is done outside the transaction as Casbin has its own transaction handling
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

func (s *UserService) Login(loginDTO *dto.LoginUserDTO) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(loginDTO.Email)
	if err != nil {
		return nil, errors.ErrUserInvalidCredentials
	}

	if user.Password == nil {
		return nil, errors.ErrUserInvalidCredentials
	}

	if !utils.ComparePassword(*user.Password, loginDTO.Password) {
		return nil, errors.ErrUserInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetUserWithPermissions(userID uint) (*dto.UserWithPermissionsDTO, error) {
	// Get user from database
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// Get user's roles from database
	_, roles, err := s.userRepo.GetUserWithRoles(userID)
	if err != nil {
		return nil, err
	}

	// Convert roles to DTO
	roleDTOs := make([]dto.RoleDTO, len(roles))
	for i, role := range roles {
		roleDTOs[i] = dto.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	// Get all permissions from Casbin for each role
	permissionMap := make(map[string]bool) // Use map to avoid duplicates
	for _, role := range roles {
		permissions, err := s.casbinService.GetPermissionsForRole(role.Name)
		if err != nil {
			// Log error but continue - don't fail completely
			continue
		}

		// permissions format: [[role, resource, action], ...]
		for _, perm := range permissions {
			if len(perm) >= 3 {
				// Create unique key: "resource:action"
				key := perm[1] + ":" + perm[2]
				permissionMap[key] = true
			}
		}
	}

	// Convert permission map to DTO array
	permissionDTOs := make([]dto.PermissionDTO, 0, len(permissionMap))
	for key := range permissionMap {
		// Split key back to resource:action
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			permissionDTOs = append(permissionDTOs, dto.PermissionDTO{
				Resource: parts[0],
				Action:   parts[1],
			})
		}
	}

	// Build response DTO
	response := &dto.UserWithPermissionsDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		AvatarURL:   user.AvatarURL,
		Bio:         user.Bio,
		Status:      user.Status,
		Roles:       roleDTOs,
		Permissions: permissionDTOs,
	}

	return response, nil
}
