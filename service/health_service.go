package service

import (
	"context"

	kithttp "github.com/go-kit/kit/transport/http"
)

var getHealthServer *kithttp.Server

func init() {
	getHealthServer = kithttp.NewServer(
		getHealth,
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)
}

type GetHealthResponse struct {
	Application string `json:"application"`
	Status      string `json:"status"`
}

func getHealth(_ context.Context, _ interface{}) (interface{}, error) {
	return &GetHealthResponse{
		Application: "maat",
		Status:      "healthy",
	}, nil
}
