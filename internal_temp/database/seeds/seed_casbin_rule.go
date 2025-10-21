package seeds

import (
	"log"

	casbinService "github.com/FeisalDy/nogo/internal/common/casbin"
	"gorm.io/gorm"
)

// SeedCasbinPolicies seeds Casbin permissions using Casbin's API
// This runs automatically on server start and updates permissions
// It will add new permissions without duplicating existing ones
func SeedCasbinPolicies(db *gorm.DB) error {
	log.Println("🌱 Auto-seeding Casbin permissions...")

	// Get Casbin service (enforcer should be initialized by main.go)
	enforcer := casbinService.GetEnforcer()
	if enforcer == nil {
		log.Println("⚠️  Casbin enforcer not initialized yet, skipping seed")
		return nil
	}

	casbin := casbinService.NewCasbinService(db)

	// ==========================================
	// ADMIN ROLE - Full Access
	// ==========================================
	adminPerms := []struct {
		resource string
		action   string
	}{
		{"users", "read"},
		{"users", "write"},
		{"users", "delete"},
		{"novels", "read"},
		{"novels", "write"},
		{"novels", "delete"},
		{"chapters", "read"},
		{"chapters", "write"},
		{"chapters", "delete"},
		{"genres", "read"},
		{"genres", "write"},
		{"genres", "delete"},
		{"tags", "read"},
		{"tags", "write"},
		{"tags", "delete"},
		{"roles", "read"},
		{"roles", "write"},
		{"media", "read"},
		{"media", "write"},
		{"media", "delete"},
	}

	for _, perm := range adminPerms {
		if err := casbin.AddPermissionForRole("admin", perm.resource, perm.action); err != nil {
			log.Printf("⚠️  Warning: Failed to add admin permission %s:%s - %v", perm.resource, perm.action, err)
		}
	}

	// ==========================================
	// AUTHOR ROLE - Content Creation
	// ==========================================
	authorPerms := []struct {
		resource string
		action   string
	}{
		{"novels", "read"},
		{"novels", "write"},
		{"chapters", "read"},
		{"chapters", "write"},
		{"genres", "read"},
		{"tags", "read"},
		{"profile", "read"},
		{"profile", "write"},
		{"media", "write"},
	}

	for _, perm := range authorPerms {
		if err := casbin.AddPermissionForRole("author", perm.resource, perm.action); err != nil {
			log.Printf("⚠️  Warning: Failed to add author permission %s:%s - %v", perm.resource, perm.action, err)
		}
	}

	// ==========================================
	// USER ROLE - Read Only
	// ==========================================
	userPerms := []struct {
		resource string
		action   string
	}{
		{"novels", "read"},
		{"chapters", "read"},
		{"genres", "read"},
		{"tags", "read"},
		{"profile", "read"},
		{"profile", "write"},
	}

	for _, perm := range userPerms {
		if err := casbin.AddPermissionForRole("user", perm.resource, perm.action); err != nil {
			log.Printf("⚠️  Warning: Failed to add user permission %s:%s - %v", perm.resource, perm.action, err)
		}
	}

	// Get final count
	allPolicies, _ := enforcer.GetPolicy()
	log.Printf("✅ Casbin auto-seed complete! Total permissions: %d", len(allPolicies))

	return nil
}
