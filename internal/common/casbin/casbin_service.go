package casbin

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/FeisalDy/nogo/internal/common/errors"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var (
	enforcer     *casbin.Enforcer
	enforcerOnce sync.Once
)

func InitCasbin(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	var err error
	enforcerOnce.Do(func() {
		adapter, adapterErr := gormadapter.NewAdapterByDB(db)
		if adapterErr != nil {
			err = fmt.Errorf("failed to create casbin adapter: %w", adapterErr)
			return
		}

		enforcer, err = casbin.NewEnforcer(modelPath, adapter)
		if err != nil {
			err = fmt.Errorf("failed to create casbin enforcer: %w", err)
			return
		}

		if loadErr := enforcer.LoadPolicy(); loadErr != nil {
			err = fmt.Errorf("failed to load policies: %w", loadErr)
			return
		}
	})

	return enforcer, err
}

func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

type CasbinService struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
}

func NewCasbinService(db *gorm.DB) *CasbinService {
	return &CasbinService{
		enforcer: GetEnforcer(),
		db:       db,
	}
}

// === Role Permission Management ===

// AddPermissionForRole adds a permission to a role
// resource: e.g., "users", "novels", "chapters"
// action: e.g., "read", "write", "delete"
func (s *CasbinService) AddPermissionForRole(roleName, resource, action string) error {
	_, err := s.enforcer.AddPolicy(roleName, resource, action)
	if err != nil {
		return errors.ErrCasbinPolicySaveFailed
	}
	return s.enforcer.SavePolicy()
}

// RemovePermissionForRole removes a permission from a role
func (s *CasbinService) RemovePermissionForRole(roleName, resource, action string) error {
	_, err := s.enforcer.RemovePolicy(roleName, resource, action)
	if err != nil {
		return errors.ErrCasbinPolicyRemoveFailed
	}
	return s.enforcer.SavePolicy()
}

func (s *CasbinService) GetPermissionsForRole(roleName string) ([][]string, error) {
	return s.enforcer.GetFilteredPolicy(0, roleName)
}

// === User Role Assignment ===

// AssignRoleToUser assigns a role to a user
func (s *CasbinService) AssignRoleToUser(userID uint, roleName string) error {
	userSubject := fmt.Sprintf("user:%d", userID)
	_, err := s.enforcer.AddRoleForUser(userSubject, roleName)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// RemoveRoleFromUser removes a role from a user
func (s *CasbinService) RemoveRoleFromUser(userID uint, roleName string) error {
	userSubject := fmt.Sprintf("user:%d", userID)
	_, err := s.enforcer.DeleteRoleForUser(userSubject, roleName)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// GetRolesForUser returns all roles for a user
func (s *CasbinService) GetRolesForUser(userID uint) ([]string, error) {
	userSubject := fmt.Sprintf("user:%d", userID)
	return s.enforcer.GetRolesForUser(userSubject)
}

// GetUsersForRole returns all user IDs that have a specific role
func (s *CasbinService) GetUsersForRole(roleName string) ([]uint, error) {
	users, err := s.enforcer.GetUsersForRole(roleName)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uint, 0, len(users))
	for _, user := range users {
		// Parse "user:123" format
		var userID uint64
		_, err := fmt.Sscanf(user, "user:%d", &userID)
		if err == nil {
			userIDs = append(userIDs, uint(userID))
		}
	}
	return userIDs, nil
}

// === Permission Checking ===

// Enforce checks if a user has permission to perform an action on a resource
func (s *CasbinService) Enforce(userID uint, resource, action string) (bool, error) {
	userSubject := fmt.Sprintf("user:%d", userID)
	return s.enforcer.Enforce(userSubject, resource, action)
}

// === Batch Operations ===

// AddPermissionsForRole adds multiple permissions to a role at once
func (s *CasbinService) AddPermissionsForRole(roleName string, permissions [][]string) error {
	rules := make([][]string, len(permissions))
	for i, perm := range permissions {
		if len(perm) != 2 {
			return fmt.Errorf("invalid permission format, expected [resource, action]")
		}
		rules[i] = []string{roleName, perm[0], perm[1]}
	}

	_, err := s.enforcer.AddPolicies(rules)
	if err != nil {
		return fmt.Errorf("failed to add permissions: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// RemoveAllPermissionsForRole removes all permissions for a role
func (s *CasbinService) RemoveAllPermissionsForRole(roleName string) error {
	_, err := s.enforcer.RemoveFilteredPolicy(0, roleName)
	if err != nil {
		return fmt.Errorf("failed to remove permissions: %w", err)
	}
	return s.enforcer.SavePolicy()
}

// === Role Management ===

// DeleteRole removes a role and all its assignments
func (s *CasbinService) DeleteRole(roleName string) error {
	// Remove all permissions for the role
	if err := s.RemoveAllPermissionsForRole(roleName); err != nil {
		return err
	}

	// Remove all user-role assignments
	_, err := s.enforcer.RemoveFilteredGroupingPolicy(1, roleName)
	if err != nil {
		return fmt.Errorf("failed to remove role assignments: %w", err)
	}

	return s.enforcer.SavePolicy()
}

// UpdateRoleName updates a role name (requires removing and re-adding policies)
func (s *CasbinService) UpdateRoleName(oldRoleName, newRoleName string) error {
	// Get all permissions for the old role
	permissions, err := s.GetPermissionsForRole(oldRoleName)
	if err != nil {
		return fmt.Errorf("failed to get permissions for role: %w", err)
	}

	// Get all users with the old role
	users, err := s.enforcer.GetUsersForRole(oldRoleName)
	if err != nil {
		return fmt.Errorf("failed to get users for role: %w", err)
	}

	// Remove the old role
	if err := s.DeleteRole(oldRoleName); err != nil {
		return err
	}

	// Add permissions with new role name
	for _, perm := range permissions {
		if len(perm) >= 3 {
			if err := s.AddPermissionForRole(newRoleName, perm[1], perm[2]); err != nil {
				return err
			}
		}
	}

	// Reassign users to new role
	for _, user := range users {
		// Extract user ID from "user:123" format
		var userID uint64
		_, err := fmt.Sscanf(user, "user:%d", &userID)
		if err == nil {
			if err := s.AssignRoleToUser(uint(userID), newRoleName); err != nil {
				return err
			}
		}
	}

	return s.enforcer.SavePolicy()
}

// === Utility Functions ===

// GetAllRoles returns all unique roles in the system
func (s *CasbinService) GetAllRoles() ([]string, error) {
	allPolicies, err := s.enforcer.GetPolicy()
	if err != nil {
		return nil, err
	}
	roleSet := make(map[string]bool)

	for _, policy := range allPolicies {
		if len(policy) > 0 {
			roleSet[policy[0]] = true
		}
	}

	roles := make([]string, 0, len(roleSet))
	for role := range roleSet {
		roles = append(roles, role)
	}
	return roles, nil
}

// HasRole checks if a user has a specific role
func (s *CasbinService) HasRole(userID uint, roleName string) (bool, error) {
	userSubject := fmt.Sprintf("user:%d", userID)
	return s.enforcer.HasRoleForUser(userSubject, roleName)
}

// GetAllSubjects returns all subjects (users) in the system
func (s *CasbinService) GetAllSubjects() []string {
	subjects, _ := s.enforcer.GetAllSubjects()
	return subjects
}

// ReloadPolicies reloads all policies from the database
func (s *CasbinService) ReloadPolicies() error {
	return s.enforcer.LoadPolicy()
}

// ClearAllPolicies removes all policies (use with caution!)
func (s *CasbinService) ClearAllPolicies() error {
	s.enforcer.ClearPolicy()
	return s.enforcer.SavePolicy()
}

// === Sync with Database ===

// SyncRolePermissions synchronizes Casbin policies with role_permissions table
// This ensures consistency when roles are modified through the database
func (s *CasbinService) SyncRolePermissions() error {
	// This function can be called periodically or after database changes
	// to ensure Casbin is in sync with the database
	return s.enforcer.LoadPolicy()
}

// GetUserIDFromSubject extracts user ID from subject string "user:123"
func GetUserIDFromSubject(subject string) (uint, error) {
	var userID uint64
	_, err := fmt.Sscanf(subject, "user:%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("invalid subject format: %s", subject)
	}
	return uint(userID), nil
}

// FormatUserSubject formats user ID to subject string "user:123"
func FormatUserSubject(userID uint) string {
	return "user:" + strconv.FormatUint(uint64(userID), 10)
}
