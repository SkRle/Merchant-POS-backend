package tests

import (
	"Merchant-POS-backend/model"
	"Merchant-POS-backend/module/user/repositories"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to memory database: %v", err)
	}
	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewUserRepository(db)

	t.Run("Create and FindByID", func(t *testing.T) {
		user := &model.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
		}
		err := repo.Create(user)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found, err := repo.FindByID(user.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found.Username != "testuser" {
			t.Errorf("Expected testuser, got %s", found.Username)
		}

		_, err = repo.FindByID(999)
		if err == nil {
			t.Errorf("Expected error for missing record, got none")
		}
	})

	t.Run("FindByEmail", func(t *testing.T) {
		found, err := repo.FindByEmail("test@example.com")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found.Email != "test@example.com" {
			t.Errorf("Expected test@example.com, got %s", found.Email)
		}

		_, err = repo.FindByEmail("missing@example.com")
		if err == nil {
			t.Errorf("Expected error for missing email, got none")
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		repo.Create(&model.User{Username: "user2", Email: "user2@example.com", Password: "pwd"})

		users, err := repo.FindAll()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(users) < 2 {
			t.Errorf("Expected at least 2 users, got %d", len(users))
		}
	})

	t.Run("Update", func(t *testing.T) {
		user, _ := repo.FindByEmail("test@example.com")
		user.Username = "updated_user"
		err := repo.Update(user)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		updated, _ := repo.FindByID(user.ID)
		if updated.Username != "updated_user" {
			t.Errorf("Expected updated_user, got %s", updated.Username)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		user, _ := repo.FindByEmail("test@example.com")
		err := repo.Delete(user.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = repo.FindByID(user.ID)
		if err == nil {
			t.Errorf("Expected error after deletion, got none")
		}
	})

	t.Run("Errors After Drop Table", func(t *testing.T) {
		db.Migrator().DropTable(&model.User{})

		_, err := repo.FindAll()
		if err == nil {
			t.Errorf("Expected error for FindAll after dropping table")
		}
	})
}
