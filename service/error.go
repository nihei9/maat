package service

import "encoding/json"

// ErrorResponse is a response used when an error occured.
// It implements interfaces below.
// * error
// * encoding/json.Marshaler
// * github.com/go-kit/kit/transport/http.StatusCoder
type ErrorResponse struct {
	err        error
	statusCode int
}

func NewErrorResponse(err error, statusCode int) error {
	return &ErrorResponse{
		err:        err,
		statusCode: statusCode,
	}
}

// Error is an implementation of error#Error.
func (e *ErrorResponse) Error() string {
	return e.err.Error()
}

// MarshalJSON is an implementation of encoding/json.Marshaler#MarshalJSON.
func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	data := struct {
		Message string `json:"message"`
	}{
		Message: e.err.Error(),
	}

	return json.Marshal(&data)
}

// StatusCode is an implementation of github.com/go-kit/kit/transport/http.StatusCoder#StatusCode.
func (e *ErrorResponse) StatusCode() int {
	return e.statusCode
}
