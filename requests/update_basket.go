package requests

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jamess-Lucass/microservice-shop/basket-service/models"
	"github.com/Jamess-Lucass/microservice-shop/basket-service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type UpdateBasketRequest struct {
	Items []UpdateBasketItemRequest `json:"items" validate:"dive"`
}

type UpdateBasketItemRequest struct {
	Id        string `json:"id" validate:"omitempty,uuid"`
	CatalogId string `json:"catalogId" validate:"required,uuid"`
	Quantity  uint   `json:"quantity" validate:"required,min=1"`
}

func (r *UpdateBasketRequest) Bind(c *fiber.Ctx, basket *models.Basket, v *validator.Validate) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Struct(r); err != nil {
		return err
	}

	var items []models.BasketItem

	for _, item := range r.Items {
		uri := fmt.Sprintf("%s/api/v1/catalog/%s", os.Getenv("CATALOG_SERVICE_BASE_URL"), item.CatalogId)
		body, err := utils.HttpGet(uri)
		if err != nil {
			return err
		}

		var catalog CatalogResponse
		if err := json.NewDecoder(body).Decode(&catalog); err != nil {
			return err
		}

		if item.Id != "" {
			basketItem, ok := lo.Find(basket.Items, func(i models.BasketItem) bool {
				return i.ID == uuid.MustParse(item.Id)
			})

			if !ok {
				return fmt.Errorf("could not find item with the id: %s", item.Id)
			}

			if uuid.MustParse(item.CatalogId) != basketItem.CatalogId {
				return fmt.Errorf("cannot change the catalogId, you must remove the item from your basket")
			}

			basketItem.Quantity = item.Quantity

			items = append(items, basketItem)
		} else {
			basketItem := models.BasketItem{
				ID:        uuid.New(),
				CatalogId: uuid.MustParse(item.CatalogId),
				Price:     catalog.Price,
				Quantity:  item.Quantity,
			}

			items = append(items, basketItem)
		}
	}

	basket.Items = items

	return nil
}
