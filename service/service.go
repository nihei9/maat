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
