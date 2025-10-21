package repository

import (
	commonModel "github.com/FeisalDy/nogo/internal/common/model"
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
// Note: These methods only deal with the user_roles junction table (common domain)
// They work with role IDs only, not role entities (to maintain domain boundaries)

// AssignRoleToUser assigns a role to a user by creating an entry in user_roles table
func (r *UserRepository) AssignRoleToUser(userID, roleID uint) error {
	userRole := commonModel.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

// RemoveRoleFromUser removes a role from a user by deleting from user_roles table
func (r *UserRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&commonModel.UserRole{}).Error
}

// GetUserRoleIDs gets all role IDs for a user from user_roles table
func (r *UserRepository) GetUserRoleIDs(userID uint) ([]uint, error) {
	var roleIDs []uint
	err := r.db.
		Model(&commonModel.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// HasRoleByID checks if a user has a specific role by checking user_roles table
// Only checks the junction table, doesn't access the roles table
func (r *UserRepository) HasRoleByID(userID, roleID uint) (bool, error) {
	var count int64
	err := r.db.
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	return count > 0, err
}
