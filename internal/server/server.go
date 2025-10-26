package server

import (
	"github.com/gofiber/fiber/v2"

	"inventory-modular-monolith/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "inventory-modular-monolith",
			AppName:      "inventory-modular-monolith",
		}),

		db: database.New(),
	}

	return server
}
