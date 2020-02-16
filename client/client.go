package client

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/nihei9/maat/service"
	"net/url"
)

type Endpoints struct {
	GetHealth             endpoint.Endpoint
	PostValidation        endpoint.Endpoint
	PostValidationTargets endpoint.Endpoint
}

func MakeEndpoints(host string, options ...kithttp.ClientOption) (*Endpoints, error) {
	target, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	target.Path = ""

	return &Endpoints{
		GetHealth: kithttp.NewClient(
			"GET",
			target,
			service.EncodeGetHealthRequest,
			service.DecodeGetHealthResponse,
			options...,
		).Endpoint(),
		PostValidation: kithttp.NewClient(
			"POST",
			target,
			service.EncodePostValidationRequest,
			service.DecodePostValidationResponse,
			options...,
		).Endpoint(),
		PostValidationTargets: kithttp.NewClient(
			"POST",
			target,
			service.EncodePostValidationTargetsRequest,
			service.DecodePostValidationTargetsResponse,
			options...,
		).Endpoint(),
	}, nil
}
