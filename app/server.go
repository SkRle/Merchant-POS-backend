package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "Merchant-POS-backend/docs" // Required for Swagger docs
	"Merchant-POS-backend/module/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// StartServer sets up and starts the Fiber HTTP server
func (a *App) StartServer() error {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"code":    code,
				"status":  "error",
				"message": err.Error(),
			})
		},
	})

	// Middlewares
	fiberApp.Use(recover.New())
	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New())

	// Health check endpoint
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger UI
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	// API Routes Group
	api := fiberApp.Group("/api/v1")

	// Register Modules
	user.RegisterUserModule(api, a.DB)

	// Graceful Shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Gracefully shutting down...")
		_ = fiberApp.ShutdownWithTimeout(5 * time.Second)
	}()

	addr := fmt.Sprintf(":%d", a.Config.Server.Port)
	log.Printf("Server is starting on %s", addr)

	if err := fiberApp.Listen(addr); err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}
