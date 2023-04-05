package requests

import (
	"github.com/Jamess-Lucass/ecommerce-basket-service/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CheckoutBasketRequest struct {
	Address     string `json:"address" validate:"required,min=3,max=512"`
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name" validate:"required,min=3,max=128"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=8,max=512"`
}

type CreateOrderRequest struct {
	Address     string        `json:"address"`
	Email       string        `json:"email"`
	Name        string        `json:"name"`
	PhoneNumber string        `json:"phoneNumber"`
	Basket      models.Basket `json:"basket"`
}

func (r *CheckoutBasketRequest) Bind(c *fiber.Ctx, order *CreateOrderRequest, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Struct(r); err != nil {
		return err
	}

	order.Address = r.Address
	order.Email = r.Email
	order.Name = r.Name
	order.PhoneNumber = r.PhoneNumber

	return nil
}
