package seeds

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User model for seeding
type User struct {
	gorm.Model
	Username  *string `json:"username"`
	Email     string  `json:"email" gorm:"unique;not null;index"`
	Password  *string `json:"-"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio" gorm:"type:text"`
	Status    string  `json:"status" gorm:"default:'active';index"`
}

// UserRole model for seeding
type UserRole struct {
	gorm.Model
	UserID uint  `gorm:"not null;index;uniqueIndex:idx_user_role"`
	RoleID uint  `gorm:"not null;index;uniqueIndex:idx_user_role"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role   *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// SeedUsers seeds default users with roles
func SeedUsers(db *gorm.DB) error {
	log.Println("🌱 Seeding users...")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordStr := string(hashedPassword)

	users := []struct {
		User     User
		RoleName string
	}{
		{
			User: User{
				Username:  strPtr("admin"),
				Email:     "admin@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=Admin&background=random"),
				Bio:       strPtr("System administrator"),
				Status:    "active",
			},
			RoleName: "admin",
		},
		{
			User: User{
				Username:  strPtr("author1"),
				Email:     "author1@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=Author+One&background=random"),
				Bio:       strPtr("Professional novelist and translator"),
				Status:    "active",
			},
			RoleName: "author",
		},
		{
			User: User{
				Username:  strPtr("john_doe"),
				Email:     "john@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=John+Doe&background=random"),
				Bio:       strPtr("Novel enthusiast and reader"),
				Status:    "active",
			},
			RoleName: "user",
		},
	}

	for _, userData := range users {
		var existing User
		result := db.Where("email = ?", userData.User.Email).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// Create user
			if err := db.Create(&userData.User).Error; err != nil {
				log.Printf("⚠️  Failed to seed user %s: %v", userData.User.Email, err)
				return err
			}

			// Assign role
			var role Role
			if err := db.Where("name = ?", userData.RoleName).First(&role).Error; err != nil {
				log.Printf("⚠️  Role %s not found for user %s", userData.RoleName, userData.User.Email)
				continue
			}

			userRole := UserRole{
				UserID: userData.User.ID,
				RoleID: role.ID,
			}
			if err := db.Create(&userRole).Error; err != nil {
				log.Printf("⚠️  Failed to assign role to user %s: %v", userData.User.Email, err)
			}

			log.Printf("✅ Created user: %s with role: %s", userData.User.Email, userData.RoleName)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  User already exists: %s", userData.User.Email)
		}
	}

	log.Println("✅ Users seeding completed")
	return nil
}
