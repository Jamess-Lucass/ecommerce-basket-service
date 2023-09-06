package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var ctx = context.Background()

var Tracer = otel.Tracer("github.com/Jamess-Lucass/ecommerce-basket-service")

func Connect(log *zap.Logger) *redis.Client {
	server := os.Getenv("REDIS_HOST")
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Sugar().Fatalf("Could not parse PORT to an integar: %v", err)
	}
	pass := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", server, port),
		Password: pass,
		DB:       0,
	})

	// Configure OpenTelemtry to Elastic
	tracingProvider, err := apmotel.NewTracerProvider()
	if err != nil {
		log.Sugar().Fatalf("could not create tracing provider: %v", err)
	}
	otel.SetTracerProvider(tracingProvider)

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Sugar().Errorf("error instrumenting tracing: %v", err)
	}

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Sugar().Fatalf("error pinging redis database: %v", err)
	}

	return rdb
}
