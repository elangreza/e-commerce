package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type NotFound struct {
	Message string
}

func (e NotFound) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "not found"
}

// http error code
func (e *NotFound) HttpCode() int {
	return http.StatusNotFound
}

// grpc error code
func (e *NotFound) GrpcCode() codes.Code {
	return codes.NotFound
}
