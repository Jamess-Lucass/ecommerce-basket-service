package handlers

import (
	"encoding/json"

	"github.com/Jamess-Lucass/ecommerce-basket-service/models"
	"github.com/Jamess-Lucass/ecommerce-basket-service/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

// GET: /api/v1/baskets/1
func (s *Server) GetBasket(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	basket, err := s.basketService.Get(c.Context(), id)
	if err != nil {
		s.logger.Sugar().Debugf("error getting basket: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: "Could not retrieve basket"})
	}

	return c.Status(fiber.StatusOK).JSON(basket)
}

// POST: /api/v1/baskets
func (s *Server) CreateBasket(c *fiber.Ctx) error {
	var basket models.Basket

	req := &requests.CreateBasketRequest{}
	if err := req.Bind(c, &basket, s.validator); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	if err := s.basketService.Set(c.Context(), &basket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(basket)
}

// PUT: /api/v1/baskets/1
func (s *Server) UpdateBasket(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	basket, err := s.basketService.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	req := &requests.UpdateBasketRequest{}
	if err := req.Bind(c, basket, s.validator); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	if err := s.basketService.Set(c.Context(), basket); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(basket)
}

// DELETE: /api/v1/baskets/1
func (s *Server) DeleteBasket(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	if err := s.basketService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// POST: /api/v1/baskets/1/checkout
func (s *Server) CheckoutBasket(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	basket, err := s.basketService.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	order := requests.CreateOrderRequest{Basket: *basket}
	req := &requests.CheckoutBasketRequest{}
	if err := req.Bind(c, &order, s.validator); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	// Send to Rabbit MQ
	q, err := s.rabbitMQ.QueueDeclare(
		"orders", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		s.logger.Sugar().Errorf("error occured delcaring rabbitMQ queue: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	body, err := json.Marshal(order)
	if err != nil {
		s.logger.Sugar().Errorf("error occured delcaring rabbitMQ queue: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	if err := s.rabbitMQ.PublishWithContext(c.Context(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			Body: body,
		}); err != nil {
		s.logger.Sugar().Errorf("error occured publishing to rabbitMQ queue: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	if err := s.basketService.Delete(c.Context(), id); err != nil {
		s.logger.Sugar().Errorf("error occured publishing to rabbitMQ queue: %v", err)
	}

	return c.Status(fiber.StatusAccepted).JSON(order)
}
