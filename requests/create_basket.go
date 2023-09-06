package requests

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jamess-Lucass/ecommerce-basket-service/models"
	"github.com/Jamess-Lucass/ecommerce-basket-service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateBasketRequest struct {
	Items []CreateBasketItemRequest `json:"items" validate:"dive"`
}

type CreateBasketItemRequest struct {
	CatalogId string `json:"catalogId" validate:"required,uuid"`
	Quantity  uint   `json:"quantity" validate:"required,min=1"`
}

type CatalogResponse struct {
	Price float32 `json:"price"`
}

func (r *CreateBasketRequest) Bind(c *fiber.Ctx, basket *models.Basket, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Struct(r); err != nil {
		return err
	}

	basket.ID = uuid.New()

	for _, item := range r.Items {
		uri := fmt.Sprintf("%s/api/v1/catalog/%s", os.Getenv("CATALOG_SERVICE_BASE_URL"), item.CatalogId)
		body, err := utils.HttpGet(c.Context(), uri)
		if err != nil {
			return err
		}
		defer body.Close()

		var catalog CatalogResponse
		if err := json.NewDecoder(body).Decode(&catalog); err != nil {
			return err
		}

		basketItem := models.BasketItem{
			ID:        uuid.New(),
			CatalogId: uuid.MustParse(item.CatalogId),
			Price:     catalog.Price,
			Quantity:  item.Quantity,
		}

		basket.Items = append(basket.Items, basketItem)
	}

	return nil
}
