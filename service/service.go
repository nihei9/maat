package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHTTPHandler() http.Handler {
	r := mux.NewRouter()
	r.Methods("GET").Path("/health").Handler(getHealthServer)
	r.Methods("POST").Path("/validation").Handler(postValidationServer)
	r.Methods("POST").Path("/validation/targets").Handler(postValidationTargetsServer)

	return r
}
