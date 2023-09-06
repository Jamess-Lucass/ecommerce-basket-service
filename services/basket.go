package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Jamess-Lucass/ecommerce-basket-service/database"
	"github.com/Jamess-Lucass/ecommerce-basket-service/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type BasketService struct {
	db *redis.Client
}

func NewBasketService(db *redis.Client) *BasketService {
	return &BasketService{
		db: db,
	}
}

func (s *BasketService) Get(ctx context.Context, id uuid.UUID) (*models.Basket, error) {
	ctx, span := database.Tracer.Start(ctx, "redis")

	value, err := s.db.Get(ctx, id.String()).Result()
	if err != nil {
		return nil, err
	}

	span.End()

	var basket models.Basket
	if err := json.Unmarshal([]byte(value), &basket); err != nil {
		return nil, err
	}

	return &basket, nil
}

func (s *BasketService) Set(ctx context.Context, basket *models.Basket) error {
	value, err := json.Marshal(basket)
	if err != nil {
		return err
	}

	ctx, span := database.Tracer.Start(ctx, "redis")

	err = s.db.Set(ctx, basket.ID.String(), value, 24*time.Hour).Err()

	span.End()

	return err
}

func (s *BasketService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := database.Tracer.Start(ctx, "redis")

	err := s.db.Del(ctx, id.String()).Err()

	span.End()

	return err
}
