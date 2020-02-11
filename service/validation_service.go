package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/nihei9/maat/service/value"
)

var postValidationServer *kithttp.Server

func init() {
	postValidationServer = kithttp.NewServer(
		postValidation,
		decodePostValidationRequest,
		kithttp.EncodeJSONResponse,
	)
}

type PostValidationRequest struct {
	Expected map[string]value.Value
	Actual   map[string]value.Value
}

type PostValidationResponse struct {
	Passed bool `json:"passed"`
}

func postValidation(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*PostValidationRequest)

	for key, actual := range r.Actual {
		expected := r.Expected[key]
		if passed := expected.Test(actual); !passed {
			return PostValidationResponse{
				Passed: false,
			}, nil
		}
	}

	return &PostValidationResponse{
		Passed: true,
	}, nil
}

func decodePostValidationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	src := struct {
		Expected interface{} `json:"expected"`
		Actual   interface{} `json:"actual"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&src)
	if err != nil {
		return nil, NewErrorResponse(err, http.StatusBadRequest)
	}

	expected := map[string]value.Value{}
	{
		srcExpected, ok := src.Expected.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("'expected' field must be a map[string]interface{}")
			return nil, NewErrorResponse(err, http.StatusBadRequest)
		}

		for key, srcElem := range srcExpected {
			e, err := unmarshalValue(srcElem)
			if err != nil {
				return nil, NewErrorResponse(err, http.StatusBadRequest)
			}
			expected[key] = e
		}
	}

	actual := map[string]value.Value{}
	{
		srcActual, ok := src.Actual.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("'actual' field must be a map[string]interface{}")
			return nil, NewErrorResponse(err, http.StatusBadRequest)
		}

		for key, srcElem := range srcActual {
			e, err := unmarshalValue(srcElem)
			if err != nil {
				return nil, NewErrorResponse(err, http.StatusBadRequest)
			}
			actual[key] = e
		}
	}

	return &PostValidationRequest{
		Expected: expected,
		Actual:   actual,
	}, nil
}
