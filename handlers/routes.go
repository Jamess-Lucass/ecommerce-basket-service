package handlers

import (
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

func (s *Server) Start() error {
	f := fiber.New()
	f.Use(cors.New(cors.Config{AllowOrigins: "*", AllowCredentials: true, MaxAge: 0}))

	f.Use(fiberzap.New(fiberzap.Config{
		Logger: s.logger,
	}))

	f.Post("/api/v1/baskets", s.CreateBasket)
	f.Put("/api/v1/baskets/:id", s.UpdateBasket)
	f.Get("/api/v1/baskets/:id", s.GetBasket)
	f.Delete("/api/v1/baskets/:id", s.DeleteBasket)
	f.Post("/api/v1/baskets/:id/checkout", s.CheckoutBasket)

	f.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": fiber.StatusNotFound, "message": "No resource found"})
	})

	return f.Listen(":8080")
}
