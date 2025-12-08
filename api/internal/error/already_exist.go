package errs

import (
	"fmt"
	"net/http"
)

type AlreadyExist struct {
	Name string
}

func (e AlreadyExist) Error() string {
	if e.Name == "" {
		return "already exist"
	}

	return fmt.Sprintf("%s already exist", e.Name)
}

func (a AlreadyExist) HttpStatusCode() int {
	return http.StatusConflict
}
