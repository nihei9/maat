package service

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

const healthServicePath = "/health"

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

func EncodeGetHealthRequest(_ context.Context, req *http.Request, _ interface{}) error {
	req.URL.Path = healthServicePath
	return nil
}

func DecodeGetHealthResponse(_ context.Context, res *http.Response) (interface{}, error) {
	response := &GetHealthResponse{}
	err := json.NewDecoder(res.Body).Decode(response)
	return response, err
}
