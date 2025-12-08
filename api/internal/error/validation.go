package errs

import (
	"fmt"
	"net/http"
)

type ValidationError struct {
	Message string
}

func (v ValidationError) Error() string {
	if v.Message != "" {
		return fmt.Sprintf("validation error: %s", v.Message)
	}

	return "validation error"
}

func (a ValidationError) HttpStatusCode() int {
	return http.StatusBadRequest
}
