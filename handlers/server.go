package handlers

import (
	"github.com/Jamess-Lucass/ecommerce-basket-service/services"
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Server struct {
	validator     *validator.Validate
	logger        *zap.Logger
	rabbitMQ      *amqp091.Channel
	basketService *services.BasketService
}

func NewServer(logger *zap.Logger, rabbitMQ *amqp091.Channel, basketService *services.BasketService) *Server {
	return &Server{
		validator:     validator.New(),
		logger:        logger,
		rabbitMQ:      rabbitMQ,
		basketService: basketService,
	}
}
