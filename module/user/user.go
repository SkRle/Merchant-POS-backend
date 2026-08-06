package user

import (
	"Merchant-POS-backend/module/user/handler"
	"Merchant-POS-backend/module/user/repositories"
	"Merchant-POS-backend/module/user/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserModule(router fiber.Router, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	userGroup := router.Group("/users")
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Get("/", userHandler.GetAllUsers)
	userGroup.Get("/:id", userHandler.GetUser)
	userGroup.Put("/:id", userHandler.UpdateUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)
}
