package errs

import (
	"fmt"
	"net/http"
)

type MethodNotAllowedError struct {
	Method string
}

func (v MethodNotAllowedError) Error() string {
	return fmt.Sprintf("method %s not allowed", v.Method)
}

func (a MethodNotAllowedError) HttpStatusCode() int {
	return http.StatusMethodNotAllowed
}
