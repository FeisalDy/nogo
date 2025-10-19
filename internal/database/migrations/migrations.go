package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID          string
	Description string
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
}

// MigrationHistory tracks applied migrations
type MigrationHistory struct {
	ID          uint   `gorm:"primaryKey"`
	MigrationID string `gorm:"unique;not null"`
	AppliedAt   int64  `gorm:"not null"`
}

// GetAllMigrations returns all migrations in order
func GetAllMigrations() []Migration {
	return []Migration{
		Migration001CreateUsers(),
		Migration002CreateAuthTable(),
		Migration003CreateMedia(),
		Migration004CreateNovels(),
		Migration005CreateChapters(),
		Migration006CreateGenresAndTags(),
		Migration007AddNovelGenresAndTags(),
		Migration008SeedRoles(),
	}
}

// RunMigrations applies all pending migrations
func RunMigrations(db *gorm.DB) error {
	// Create migration history table
	if err := db.AutoMigrate(&MigrationHistory{}); err != nil {
		return fmt.Errorf("failed to create migration history table: %v", err)
	}

	migrations := GetAllMigrations()

	for _, migration := range migrations {
		// Check if migration already applied
		var count int64
		db.Model(&MigrationHistory{}).Where("migration_id = ?", migration.ID).Count(&count)

		if count > 0 {
			fmt.Printf("Migration %s already applied, skipping\n", migration.ID)
			continue
		}

		fmt.Printf("Running migration %s: %s\n", migration.ID, migration.Description)

		// Run migration
		if err := migration.Up(db); err != nil {
			return fmt.Errorf("failed to run migration %s: %v", migration.ID, err)
		}

		// Record migration as applied
		history := MigrationHistory{
			MigrationID: migration.ID,
			AppliedAt:   time.Now().Unix(), // Current timestamp
		}
		if err := db.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %v", migration.ID, err)
		}

		fmt.Printf("Migration %s completed successfully\n", migration.ID)
	}

	return nil
}
