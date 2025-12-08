package errs

import (
	"net/http"
)

type InvalidCredential struct{}

func (e InvalidCredential) Error() string {
	return "invalid credential"
}

func (a InvalidCredential) HttpStatusCode() int {
	return http.StatusUnauthorized
}
