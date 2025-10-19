package migrations

import (
	"gorm.io/gorm"
)

// ===========================
// AUTH DOMAIN MIGRATION 002
// ===========================
// NOTE: Permissions are managed by Casbin through casbin_rule table
//       No need for separate permission tables

// Role table defines user roles (e.g., admin, author, reader)
// Roles are synced with Casbin for permission management
type Role struct {
	gorm.Model
	Name        string  `gorm:"unique;not null;index"`
	Description *string `gorm:"type:text"`
}

// UserRole is a many-to-many relation between users and roles
// This is synced with Casbin's role assignments
type UserRole struct {
	gorm.Model
	UserID uint  `gorm:"not null;index;uniqueIndex:idx_user_role"`
	RoleID uint  `gorm:"not null;index;uniqueIndex:idx_user_role"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role   *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Migration002CreateAuthTable creates roles and user_roles tables
// Casbin automatically creates casbin_rule table for permissions
func Migration002CreateAuthTable() Migration {
	return Migration{
		ID:          "002_create_auth_table",
		Description: "Create auth-related tables: roles and user_roles (Casbin handles permissions)",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(
				&Role{},
				&UserRole{},
			)
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				&UserRole{},
				&Role{},
			)
		},
	}
}
