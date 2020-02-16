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

const validationTargetsServicePath = "/validation/targets"

var postValidationTargetsServer *kithttp.Server

func init() {
	postValidationTargetsServer = kithttp.NewServer(
		postValidationTargets,
		decodePostValidationTargetsRequest,
		kithttp.EncodeJSONResponse,
	)
}

type PostValidationTargetsRequest struct {
	ValidationID validation.ID          `json:"validation_id"`
	Actual       map[string]value.Value `json:"actual"`
}

func (r *PostValidationTargetsRequest) Validate() error {
	if r.ValidationID.IsNil() {
		return fmt.Errorf("'validation_id' field is required")
	}

	if len(r.Actual) <= 0 {
		return fmt.Errorf("'actual' field must contain at least one element")
	}

	return nil
}

type PostValidationTargetsResponse struct {
	httpResponse
	Passed bool `json:"passed"`
}

func postValidationTargets(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*PostValidationTargetsRequest)

	v, err := validation.Store.Load(r.ValidationID)
	if err != nil {
		return nil, NewErrorResponse(err, http.StatusInternalServerError)
	}
	if v == nil {
		err := fmt.Errorf("'validation_id' field is invalid")
		return nil, NewErrorResponse(err, http.StatusBadRequest)
	}

	for name, actual := range r.Actual {
		passed, err := v.Do(name, actual)
		if err != nil {
			return nil, NewErrorResponse(err, http.StatusInternalServerError)
		}
		if !passed {
			return PostValidationTargetsResponse{
				Passed: false,
			}, nil
		}
	}

	return &PostValidationTargetsResponse{
		Passed: true,
	}, nil
}

func decodePostValidationTargetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	src := struct {
		ValidationID validation.ID `json:"validation_id"`
		Actual       interface{}   `json:"actual"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&src)
	if err != nil {
		return nil, NewErrorResponse(err, http.StatusBadRequest)
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

	res := &PostValidationTargetsRequest{
		ValidationID: src.ValidationID,
		Actual:       actual,
	}
	err = res.Validate()
	if err != nil {
		return nil, NewErrorResponse(err, http.StatusBadRequest)
	}

	return res, nil
}

func EncodePostValidationTargetsRequest(ctx context.Context, r *http.Request, request interface{}) error {
	r.URL.Path = validationTargetsServicePath
	return kithttp.EncodeJSONRequest(ctx, r, request)
}

func DecodePostValidationTargetsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	res := &PostValidationTargetsResponse{}
	err := json.NewDecoder(r.Body).Decode(res)
	res.setHTTPResponse(r)
	return res, err
}
