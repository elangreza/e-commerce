package errs

import (
	"fmt"
	"net/http"
	"strings"
)

type NotFound struct {
	Message string
}

func (e NotFound) Error() string {
	if e.Message == "" {
		return "not found"
	}

	if strings.Contains(e.Message, "not found") {
		return e.Message
	}

	return fmt.Sprintf("%s not found", e.Message)
}

func (a NotFound) HttpStatusCode() int {
	return http.StatusNotFound
}
