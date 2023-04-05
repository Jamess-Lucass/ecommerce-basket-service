package services

import (
	"encoding/json"
	"time"

	"github.com/Jamess-Lucass/ecommerce-basket-service/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
)

type BasketService struct {
	db *redis.Client
}

func NewBasketService(db *redis.Client) *BasketService {
	return &BasketService{
		db: db,
	}
}

func (s *BasketService) Get(ctx *fasthttp.RequestCtx, id uuid.UUID) (*models.Basket, error) {
	value, err := s.db.Get(ctx, id.String()).Result()
	if err != nil {
		return nil, err
	}

	var basket models.Basket
	if err := json.Unmarshal([]byte(value), &basket); err != nil {
		return nil, err
	}

	return &basket, nil
}

func (s *BasketService) Set(ctx *fasthttp.RequestCtx, basket *models.Basket) error {
	value, err := json.Marshal(basket)
	if err != nil {
		return err
	}

	return s.db.Set(ctx, basket.ID.String(), value, 24*time.Hour).Err()
}

func (s *BasketService) Delete(ctx *fasthttp.RequestCtx, id uuid.UUID) error {
	return s.db.Del(ctx, id.String()).Err()
}
