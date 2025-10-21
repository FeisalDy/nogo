package repository

import (
	"github.com/FeisalDy/nogo/internal/common/model"
	"github.com/FeisalDy/nogo/internal/database"
	roleModel "github.com/FeisalDy/nogo/internal/role/model"
	userModel "github.com/FeisalDy/nogo/internal/user/model"
	"gorm.io/gorm"
)

// RoleRepository handles role-related database operations
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) WithTx(tx *gorm.DB) *RoleRepository {
	return &RoleRepository{db: tx}
}

func (r *RoleRepository) Create(role *roleModel.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByID(id uint) (*roleModel.Role, error) {
	var role roleModel.Role
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *RoleRepository) GetByName(name string) (*roleModel.Role, error) {
	var role roleModel.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepository) GetAll() ([]roleModel.Role, error) {
	var roles []roleModel.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// Update updates a role
func (r *RoleRepository) Update(role *roleModel.Role) error {
	return r.db.Save(role).Error
}

// Delete deletes a role
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&roleModel.Role{}, id).Error
}

// GetRoleWithUsers gets a role and all users with that role
func (r *RoleRepository) GetRoleWithUsers(roleID uint) (*roleModel.Role, []userModel.User, error) {
	var role roleModel.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return nil, nil, err
	}

	var users []userModel.User
	err := r.db.
		Table("users").
		Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ? AND users.deleted_at IS NULL", roleID).
		Find(&users).Error

	return &role, users, err
}

// GetUserCountByRole gets the count of users with a specific role
func (r *RoleRepository) GetUserCountByRole(roleID uint) (int64, error) {
	var count int64
	err := r.db.
		Model(&model.UserRole{}).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count, err
}

// Exists checks if a role exists by ID or name
func (r *RoleRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&roleModel.Role{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// ExistsByName checks if a role exists by name
func (r *RoleRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := database.DB.Model(&roleModel.Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
