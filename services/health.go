package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type HealthService struct {
	db *redis.Client
}

func NewHealthService(db *redis.Client) *HealthService {
	return &HealthService{
		db: db,
	}
}

func (s *HealthService) Ping(ctx context.Context) error {
	_, err := s.db.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
