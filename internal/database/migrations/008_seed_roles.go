package migrations

import (
	"gorm.io/gorm"
)

func Migration008SeedRoles() Migration {
	return Migration{
		ID:          "008_seed_roles",
		Description: "Seed default roles into the database",
		Up: func(db *gorm.DB) error {

			defaultRoles := []Role{
				{Name: "admin", Description: strPtr("Administrator with full access")},
				{Name: "editor", Description: strPtr("Can create and edit content")},
				{Name: "user", Description: strPtr("Regular user with basic access")},
			}

			for _, role := range defaultRoles {
				var existingRole Role
				if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
					if err := db.Create(&role).Error; err != nil {
						return err
					}
				}
			}
			return nil
		},
		Down: func(db *gorm.DB) error {
			roleNames := []string{"admin", "editor", "user"}

			if err := db.Where("name IN ?", roleNames).Delete(&Role{}).Error; err != nil {
				return err
			}

			return nil
		},
	}
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}
