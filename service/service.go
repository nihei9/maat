package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHTTPHandler() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path(healthServicePath).Handler(getHealthServer)
	r.Methods("POST").Path(validationServicePath).Handler(postValidationServer)
	r.Methods("POST").Path(validationTargetsServicePath).Handler(postValidationTargetsServer)

	return r
}

type httpResponse struct {
	status     string
	statusCode int
	header     http.Header
}

func (r *httpResponse) setHTTPResponse(res *http.Response) {
	r.status = res.Status
	r.statusCode = res.StatusCode
	r.header = res.Header
}

func (r *httpResponse) HTTPStatus() string {
	return r.status
}

func (r *httpResponse) HTTPStatusCode() int {
	return r.statusCode
}

func (r *httpResponse) HTTPHeader() http.Header {
	return r.header
}
