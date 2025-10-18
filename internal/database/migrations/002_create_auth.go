package migrations

import (
	"gorm.io/gorm"
)

// ===========================
// AUTH DOMAIN MIGRATION 002
// ===========================

// Role table defines user roles (e.g., admin, author, reader)
type Role struct {
	gorm.Model
	Name        string  `gorm:"unique;not null"`
	Description *string `gorm:"type:text"`
}

// Permission table defines fine-grained permissions
type Permission struct {
	gorm.Model
	Name        string  `gorm:"unique;not null"`
	Description *string `gorm:"type:text"`
}

// UserRole is a many-to-many relation between users and roles
type UserRole struct {
	gorm.Model
	UserID uint  `gorm:"not null;index;uniqueIndex:idx_user_role"`
	RoleID uint  `gorm:"not null;index;uniqueIndex:idx_user_role"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role   *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// RolePermission links roles to their permissions
type RolePermission struct {
	gorm.Model
	RoleID       uint        `gorm:"not null;index;uniqueIndex:idx_role_permission"`
	PermissionID uint        `gorm:"not null;index;uniqueIndex:idx_role_permission"`
	Role         *Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission   *Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Migration002CreateAuthTable creates roles, permissions, and relations
func Migration002CreateAuthTable() Migration {
	return Migration{
		ID:          "002_create_auth_table",
		Description: "Create auth-related tables: roles, permissions, user_roles, role_permissions",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(
				&Role{},
				&Permission{},
				&UserRole{},
				&RolePermission{},
			)
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				&RolePermission{},
				&UserRole{},
				&Permission{},
				&Role{},
			)
		},
	}
}
