package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func HttpGet(ctx context.Context, uri string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := tracingClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= fiber.StatusBadRequest {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(body))
	}

	return res.Body, nil
}
