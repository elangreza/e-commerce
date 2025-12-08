package errs

import (
	"fmt"
	"net/http"
)

type NotFound struct {
	Message string
}

func (e NotFound) Error() string {
	if e.Message == "" {
		return "not found"
	}

	return fmt.Sprintf("%s not found", e.Message)
}

func (a NotFound) HttpStatusCode() int {
	return http.StatusNotFound
}
