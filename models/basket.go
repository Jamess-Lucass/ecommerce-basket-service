package models

import (
	"github.com/google/uuid"
)

type Basket struct {
	ID    uuid.UUID    `json:"id"`
	Items []BasketItem `json:"items"`
}

type BasketItem struct {
	ID        uuid.UUID `json:"id"`
	CatalogId uuid.UUID `json:"catalogId"`
	Price     float32   `json:"price"`
	Quantity  uint      `json:"quantity"`
}
