package seeds

import (
	"log"

	"gorm.io/gorm"
)

// Role model for seeding
type Role struct {
	gorm.Model
	Name        string  `gorm:"unique;not null;index"`
	Description *string `gorm:"type:text"`
}

// SeedRoles seeds default roles
func SeedRoles(db *gorm.DB) error {
	log.Println("🌱 Seeding roles...")

	roles := []Role{
		{Name: "admin", Description: strPtr("Administrator with full access to all resources")},
		{Name: "author", Description: strPtr("Content creator who can write and manage novels and chapters")},
		{Name: "user", Description: strPtr("Regular user with read access")},
	}

	for _, role := range roles {
		var existing Role
		result := db.Where("name = ?", role.Name).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("⚠️  Failed to seed role %s: %v", role.Name, err)
				return err
			}
			log.Printf("✅ Created role: %s", role.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  Role already exists: %s", role.Name)
		}
	}

	log.Println("✅ Roles seeding completed")
	return nil
}
