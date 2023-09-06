package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Jamess-Lucass/ecommerce-basket-service/database"
	"github.com/Jamess-Lucass/ecommerce-basket-service/handlers"
	"github.com/Jamess-Lucass/ecommerce-basket-service/services"
	"github.com/rabbitmq/amqp091-go"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LOG_LEVEL = os.Getenv("LOG_LEVEL")
var LOG_LEVELS = map[string]zapcore.Level{
	"DEBUG": zap.DebugLevel,
	"INFO":  zap.InfoLevel,
	"WARN":  zap.WarnLevel,
	"ERROR": zap.ErrorLevel,
}

func main() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, LOG_LEVELS[LOG_LEVEL])
	logger := zap.New(core, zap.AddCaller())

	db := database.Connect(logger)

	// Rabbit MQ
	user := os.Getenv("RABBITMQ_USERNAME")
	pass := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	u := &url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
	}

	rabbitMQClient, err := amqp091.Dial(u.String())
	if err != nil {
		logger.Sugar().Fatalf("error occured connecting to rabbit MQ: %v", err)
	}
	defer rabbitMQClient.Close()

	ch, err := rabbitMQClient.Channel()
	if err != nil {
		logger.Sugar().Fatalf("error occured opening rabbitMQ channel: %v", err)
	}
	defer ch.Close()

	http.DefaultTransport = apmhttp.WrapRoundTripper(http.DefaultTransport, apmhttp.WithClientTrace())

	healthService := services.NewHealthService(db)
	basketService := services.NewBasketService(db)

	server := handlers.NewServer(logger, ch, healthService, basketService)

	if err := server.Start(); err != nil {
		logger.Sugar().Fatalf("error starting web server: %v", err)
	}
}
