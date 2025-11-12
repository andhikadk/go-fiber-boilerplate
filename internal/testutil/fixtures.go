package testutil

import (
	"go-fiber-boilerplate/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUserFixture creates a test user with the given parameters
func CreateUserFixture(db *gorm.DB, name, email, password, role string) *models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: true,
	}

	db.Create(user)
	return user
}

// CreateBookFixture creates a test book with the given parameters
func CreateBookFixture(db *gorm.DB, title, author, isbn string, year int) *models.Book {
	book := &models.Book{
		Title:  title,
		Author: author,
		Year:   year,
		ISBN:   isbn,
	}

	db.Create(book)
	return book
}

// CreateMultipleUserFixtures creates multiple test users for batch testing
func CreateMultipleUserFixtures(db *gorm.DB, count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateUserFixture(
			db,
			"Test User "+string(rune('A'+i)),
			"user"+string(rune('a'+i))+"@test.com",
			"password123",
			"user",
		)
	}
	return users
}

// CreateMultipleBookFixtures creates multiple test books for batch testing
func CreateMultipleBookFixtures(db *gorm.DB, count int) []*models.Book {
	books := make([]*models.Book, count)
	for i := 0; i < count; i++ {
		books[i] = CreateBookFixture(
			db,
			"Test Book "+string(rune('A'+i)),
			"Author "+string(rune('A'+i)),
			"ISBN-"+string(rune('0'+i))+"00000000",
			2020+i,
		)
	}
	return books
}

// CreateAdminUserFixture creates a test admin user
func CreateAdminUserFixture(db *gorm.DB) *models.User {
	return CreateUserFixture(db, "Admin User", "admin@test.com", "admin123", "admin")
}

// CreateStandardUserFixture creates a test standard user
func CreateStandardUserFixture(db *gorm.DB) *models.User {
	return CreateUserFixture(db, "Standard User", "user@test.com", "user123", "user")
}
