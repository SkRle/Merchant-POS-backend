package tests

import (
	"testing"

	"Merchant-POS-backend/module/user"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterUserModule(t *testing.T) {
	app := fiber.New()
	db := setupTestDB(t)

	// Simply call it to hit coverage. We know it works if it doesn't panic.
	user.RegisterUserModule(app, db)

	// Check if routes are registered
	if len(app.GetRoutes()) == 0 {
		t.Errorf("Expected routes to be registered")
	}
}
