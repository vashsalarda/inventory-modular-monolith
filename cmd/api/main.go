package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"inventory-modular-monolith/internal/config"
	"inventory-modular-monolith/internal/database"
	"inventory-modular-monolith/internal/modules/inventory"
	"inventory-modular-monolith/internal/server"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	server := server.New()

	server.RegisterFiberRoutes()

	cfg := config.Load()
	
	// Initialize MongoDB
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Disconnect()

	app := fiber.New(fiber.Config{
		AppName: "Inventory API v1.0",
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "timestamp": time.Now()})
	})

	// API v1 group
	v1 := app.Group("/api/v1")

	// Register modules
	inventory.RegisterRoutes(v1, db)
	// pos.RegisterRoutes(v1, db)
	// merchant.RegisterRoutes(v1, db)

	log.Println("Current Time:", time.Now().Format("Jan 02, 2006 03:04 AM"))

	// List all routes
	log.Println("API Routes:")

	routes := app.GetRoutes()
	for _, route := range routes {
		method := route.Method
		path := route.Path
		if path != "/" && method != "HEAD" {
			log.Printf("%-6s %-20s\n", method, path)
		}
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	<-quit
	log.Println("Shutting down server...")
	
	if err := app.ShutdownWithContext(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
