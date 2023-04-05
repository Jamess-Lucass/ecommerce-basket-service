package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func Connect(log *zap.Logger) *redis.Client {
	server := os.Getenv("REDIS_HOST")
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Sugar().Fatalf("Could not parse PORT to an integar: %v", err)
	}
	pass := os.Getenv("REDIS_PASSWORD")

	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", server, port),
		Password: pass,
		DB:       0,
	})

	if _, err := db.Ping(ctx).Result(); err != nil {
		log.Sugar().Fatalf("error pinging redis database: %v", err)
	}

	return db
}
