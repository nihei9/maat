package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/nihei9/maat/service/validation"
	"github.com/nihei9/maat/service/value"
)

const validationServicePath = "/validation"

var postValidationServer *kithttp.Server

func init() {
	postValidationServer = kithttp.NewServer(
		postValidation,
		decodePostValidationRequest,
		kithttp.EncodeJSONResponse,
	)
}

type PostValidationRequest struct {
	Expected map[string]value.Value `json:"expected"`
}

func (r *PostValidationRequest) Validate() error {
	if len(r.Expected) <= 0 {
		return fmt.Errorf("'expected' field must contain at least one element")
	}

	return nil
}

type PostValidationResponse struct {
	httpResponse
	ValidationID validation.ID `json:"validation_id"`
}

func postValidation(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*PostValidationRequest)

	v := validation.NewValidation()
	for name, expected := range r.Expected {
		v.Expect(name, expected)
	}

	id, err := validation.Store.Store(v)
	if err != nil {
		return nil, NewErrorResponse(err, http.StatusInternalServerError)
	}

	return &PostValidationResponse{
		ValidationID: id,
	}, nil
}

func decodePostValidationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	src := struct {
		Expected interface{} `json:"expected"`
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

	res := &PostValidationRequest{
		Expected: expected,
	}
	err = res.Validate()
	if err != nil {
		return nil, NewErrorResponse(err, http.StatusBadRequest)
	}

	return res, nil
}

func EncodePostValidationRequest(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = validationServicePath
	return kithttp.EncodeJSONRequest(ctx, r, request)
}

func DecodePostValidationResponse(_ context.Context, r *http.Response) (interface{}, error) {
	res := &PostValidationResponse{}
	err := json.NewDecoder(r.Body).Decode(res)
	res.setHTTPResponse(r)
	return res, err
}
