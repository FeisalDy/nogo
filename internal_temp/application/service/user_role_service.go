package service

import (
	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/FeisalDy/nogo/internal/database"
	roleRepo "github.com/FeisalDy/nogo/internal/role/repository"
	userRepo "github.com/FeisalDy/nogo/internal/user/repository"
	"gorm.io/gorm"
)

// UserRoleService handles cross-domain operations between User and Role domains
// This is part of the Application Layer that coordinates between multiple domains
type UserRoleService struct {
	userRepo      *userRepo.UserRepository
	roleRepo      *roleRepo.RoleRepository
	casbinService *casbinService.CasbinService
}

// NewUserRoleService creates a new instance of UserRoleService
func NewUserRoleService(
	userRepository *userRepo.UserRepository,
	roleRepository *roleRepo.RoleRepository,
	casbin *casbinService.CasbinService,
) *UserRoleService {
	return &UserRoleService{
		userRepo:      userRepository,
		roleRepo:      roleRepository,
		casbinService: casbin,
	}
}

// AssignRoleToUser assigns a role to a user
// This is a cross-domain operation that:
// 1. Validates user exists (User domain)
// 2. Validates role exists (Role domain)
// 3. Checks if user already has the role
// 4. Creates the user-role relationship in database
// 5. Syncs with Casbin for authorization
func (s *UserRoleService) AssignRoleToUser(userID, roleID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validate user exists
		user, err := s.userRepo.WithTx(tx).GetUserByID(userID)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.ErrUserNotFound
		}

		// 2. Validate role exists
		role, err := s.roleRepo.WithTx(tx).GetByID(roleID)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.ErrRoleNotFound
		}

		// 3. Check if user already has the role
		hasRole, err := s.userRepo.WithTx(tx).HasRoleByID(userID, roleID)
		if err != nil {
			return err
		}
		if hasRole {
			return errors.ErrUserAlreadyHasRole
		}

		// 4. Assign role in database (user_roles table)
		if err := s.userRepo.WithTx(tx).AssignRoleToUser(userID, roleID); err != nil {
			return err
		}

		// 5. Assign role in Casbin (for authorization checks)
		// Note: Casbin operations are done outside transaction as it has its own handling
		if err := s.casbinService.AssignRoleToUser(userID, role.Name); err != nil {
			return err
		}

		return nil
	})
}

// RemoveRoleFromUser removes a role from a user
// This is a cross-domain operation that:
// 1. Validates user exists (User domain)
// 2. Validates role exists (Role domain)
// 3. Checks if user has the role
// 4. Removes the user-role relationship from database
// 5. Syncs with Casbin for authorization
func (s *UserRoleService) RemoveRoleFromUser(userID, roleID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validate user exists
		user, err := s.userRepo.WithTx(tx).GetUserByID(userID)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.ErrUserNotFound
		}

		// 2. Validate role exists
		role, err := s.roleRepo.WithTx(tx).GetByID(roleID)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.ErrRoleNotFound
		}

		// 3. Check if user has the role
		hasRole, err := s.userRepo.WithTx(tx).HasRoleByID(userID, roleID)
		if err != nil {
			return err
		}
		if !hasRole {
			return errors.ErrUserDoesNotHaveRole
		}

		// 4. Remove role from database (user_roles table)
		if err := s.userRepo.WithTx(tx).RemoveRoleFromUser(userID, roleID); err != nil {
			return err
		}

		// 5. Remove role from Casbin (for authorization checks)
		if err := s.casbinService.RemoveRoleFromUser(userID, role.Name); err != nil {
			return err
		}

		return nil
	})
}

// AssignDefaultRoleToNewUser assigns the default "user" role to a newly registered user
// This is called during user registration to ensure all users have a default role
func (s *UserRoleService) AssignDefaultRoleToNewUser(tx *gorm.DB, userID uint) error {
	// Get the default "user" role
	defaultRole, err := s.roleRepo.WithTx(tx).GetByName("user")
	if err != nil {
		return err
	}
	if defaultRole == nil {
		return errors.ErrRoleNotFound
	}

	// Assign role in database (user_roles table)
	if err := s.userRepo.WithTx(tx).AssignRoleToUser(userID, defaultRole.ID); err != nil {
		return err
	}

	// Assign role in Casbin (for authorization checks)
	// Note: This is done outside the transaction as Casbin has its own transaction handling
	if err := s.casbinService.AssignRoleToUser(userID, defaultRole.Name); err != nil {
		return err
	}

	return nil
}

// GetUserRoles retrieves all roles assigned to a user
// This is a cross-domain operation that:
// 1. Validates user exists (User domain)
// 2. Gets role IDs from user_roles table (via User repository)
// 3. Fetches role details from Role domain
func (s *UserRoleService) GetUserRoles(userID uint) ([]interface{}, error) {
	// 1. Validate user exists (User domain)
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 2. Get role IDs from user_roles table (stays in User domain boundary)
	roleIDs, err := s.userRepo.GetUserRoleIDs(userID)
	if err != nil {
		return nil, err
	}

	// 3. Fetch role details from Role domain
	result := make([]interface{}, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		role, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			// Skip roles that no longer exist or have errors
			continue
		}
		if role != nil {
			result = append(result, map[string]interface{}{
				"id":          role.ID,
				"name":        role.Name,
				"description": role.Description,
			})
		}
	}

	return result, nil
}
