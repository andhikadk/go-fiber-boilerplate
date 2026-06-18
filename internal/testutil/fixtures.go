package testutil

import (
	"fmt"

	"go-fiber-boilerplate/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserFixture(db *gorm.DB, firstName, email, password, role string) *models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordStr := string(hashedPassword)
	roleStr := role
	user := &models.User{
		Email:               email,
		Password:            &passwordStr,
		PasswordIsSetByUser: true,
		Role:                &roleStr,
		IsActive:            true,
	}
	db.Create(user)
	db.Create(&models.UserProfile{UserID: user.ID, FirstName: firstName})
	return user
}

func CreateResourceFixture(db *gorm.DB, userID uint, name string) *models.Resource {
	resource := &models.Resource{
		Name:        name,
		Status:      "active",
		CreatedByID: userID,
	}
	db.Create(resource)
	return resource
}

func CreateMultipleUserFixtures(db *gorm.DB, count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateUserFixture(db, fmt.Sprintf("User %d", i+1), fmt.Sprintf("user%d@test.com", i+1), "password123", "user")
	}
	return users
}

func CreateMultipleResourceFixtures(db *gorm.DB, userID uint, count int) []*models.Resource {
	resources := make([]*models.Resource, count)
	for i := 0; i < count; i++ {
		resources[i] = CreateResourceFixture(db, userID, fmt.Sprintf("Resource %d", i+1))
	}
	return resources
}

func CreateAdminUserFixture(db *gorm.DB) *models.User {
	return CreateUserFixture(db, "Admin", "admin@test.com", "admin123", "admin")
}

func CreateStandardUserFixture(db *gorm.DB) *models.User {
	return CreateUserFixture(db, "User", "user@test.com", "user123456", "user")
}
