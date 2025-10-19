package repository

import (
	commonModel "github.com/FeisalDy/nogo/internal/common/model"
	roleModel "github.com/FeisalDy/nogo/internal/role/model"
	"github.com/FeisalDy/nogo/internal/user/model"
	"gorm.io/gorm"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID gets a user by numeric ID from the database
func (r *UserRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail gets a user by email from the database
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ===== Role-related methods =====

// GetUserWithRoles gets a user and their roles
func (r *UserRepository) GetUserWithRoles(userID uint) (*model.User, []roleModel.Role, error) {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, nil, err
	}

	var roles []roleModel.Role
	err := r.db.
		Table("roles").
		Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.deleted_at IS NULL", userID).
		Find(&roles).Error

	return &user, roles, err
}

// AssignRoleToUser assigns a role to a user
func (r *UserRepository) AssignRoleToUser(userID, roleID uint) error {
	userRole := commonModel.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

// RemoveRoleFromUser removes a role from a user
func (r *UserRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&commonModel.UserRole{}).Error
}

// GetUserRoleIDs gets all role IDs for a user
func (r *UserRepository) GetUserRoleIDs(userID uint) ([]uint, error) {
	var roleIDs []uint
	err := r.db.
		Model(&commonModel.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// HasRole checks if a user has a specific role
func (r *UserRepository) HasRole(userID uint, roleName string) (bool, error) {
	var count int64
	err := r.db.
		Table("user_roles").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ? AND roles.deleted_at IS NULL", userID, roleName).
		Count(&count).Error
	return count > 0, err
}

// HasRoleByID checks if a user has a specific role by role ID
func (r *UserRepository) HasRoleByID(userID, roleID uint) (bool, error) {
	var count int64
	err := r.db.
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	return count > 0, err
}

// HasAnyRole checks if a user has any of the specified roles
func (r *UserRepository) HasAnyRole(userID uint, roleNames []string) (bool, error) {
	var count int64
	err := r.db.
		Table("user_roles").
		Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name IN ? AND roles.deleted_at IS NULL", userID, roleNames).
		Count(&count).Error
	return count > 0, err
}
