package service

import (
	"strings"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/errors"
	roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
	userDto "github.com/FeisalDy/nogo/internal/user/dto"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
)

// UserProfileService handles user profile operations that span multiple domains
// This is part of the Application Layer
type UserProfileService struct {
	userRepo      *userRepo.UserRepository
	roleRepo      *roleRepo.RoleRepository
	casbinService *casbinService.CasbinService
}

// NewUserProfileService creates a new instance of UserProfileService
func NewUserProfileService(
	userRepository *userRepo.UserRepository,
	roleRepository *roleRepo.RoleRepository,
	casbin *casbinService.CasbinService,
) *UserProfileService {
	return &UserProfileService{
		userRepo:      userRepository,
		roleRepo:      roleRepository,
		casbinService: casbin,
	}
}

// GetUserWithPermissions retrieves user profile with roles and permissions
// This is a cross-domain operation that:
// 1. Gets user from User domain
// 2. Gets role IDs from user_roles table
// 3. Gets role details from Role domain
// 4. Gets permissions from Casbin (authorization domain)
func (s *UserProfileService) GetUserWithPermissions(userID uint) (*userDto.UserWithPermissionsDTO, error) {
	// 1. Get user from User domain
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 2. Get role IDs from user_roles table (via User repository)
	roleIDs, err := s.userRepo.GetUserRoleIDs(userID)
	if err != nil {
		return nil, err
	}

	// 3. Get role details from Role domain
	roleDTOs := make([]userDto.RoleDTO, 0, len(roleIDs))
	roleNames := make([]string, 0, len(roleIDs))

	for _, roleID := range roleIDs {
		role, err := s.roleRepo.GetByID(roleID)
		if err != nil || role == nil {
			// Skip roles that no longer exist
			continue
		}

		roleDTOs = append(roleDTOs, userDto.RoleDTO{
			ID:   role.ID,
			Name: role.Name,
		})
		roleNames = append(roleNames, role.Name)
	}

	// 4. Get permissions from Casbin for each role
	permissionMap := make(map[string]bool) // Use map to avoid duplicates
	for _, roleName := range roleNames {
		permissions, err := s.casbinService.GetPermissionsForRole(roleName)
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
	permissionDTOs := make([]userDto.PermissionDTO, 0, len(permissionMap))
	for key := range permissionMap {
		// Split key back to resource:action
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			permissionDTOs = append(permissionDTOs, userDto.PermissionDTO{
				Resource: parts[0],
				Action:   parts[1],
			})
		}
	}

	// Build response DTO
	response := &userDto.UserWithPermissionsDTO{
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
