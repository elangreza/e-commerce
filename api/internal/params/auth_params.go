package params

import (
	"strings"

	errs "github.com/elangreza/e-commerce/api/internal/error"
)

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rur *RegisterUserRequest) Validate() error {
	if rur.Name == "" {
		return errs.ValidationError{Message: "name is required"}
	}

	if !isValidEmail(rur.Email) {
		return errs.ValidationError{Message: "email is not valid"}
	}

	if rur.Password == "" {
		return errs.ValidationError{Message: "password is required"}
	}

	return nil
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lur *LoginUserRequest) Validate() error {
	if lur.Email == "" {
		return errs.ValidationError{Message: "email is required"}
	}

	if !isValidEmail(lur.Email) {
		return errs.ValidationError{Message: "email is not valid"}
	}

	if lur.Password == "" {
		return errs.ValidationError{Message: "password is required"}
	}

	return nil
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	return true
}

type ProcessTokenResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
