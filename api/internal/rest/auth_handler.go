package rest

import (
	"context"
	"encoding/json"
	"net/http"

	errs "github.com/elangreza/e-commerce/api/internal/error"
	"github.com/elangreza/e-commerce/api/internal/params"
	"github.com/go-chi/chi/v5"
)

type (
	AutService interface {
		RegisterUser(ctx context.Context, req params.RegisterUserRequest) error
		LoginUser(ctx context.Context, req params.LoginUserRequest) (string, error)
	}

	AuthHandler struct {
		svc AutService
	}
)

func NewAuthHandler(ar chi.Router, authService AutService) {

	authHandler := AuthHandler{
		svc: authService,
	}

	ar.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/login", authHandler.LoginUser)
	})
}

// RegisterUser handles user registration.
//
//	@Summary		Register User
//	@Description	Register a new user with the provided details.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		params.RegisterUserRequest	true	"Register User Request"
//	@Success		201		{string}	string						"ok"
//	@Failure		400		{object}	errs.ValidationError
//	@Failure		500		{object}	APIError
//	@Router			/auth/register [post]
func (ah *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := params.RegisterUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, errs.ValidationError{Message: err.Error()})
		return
	}

	if err := body.Validate(); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err := ah.svc.RegisterUser(r.Context(), body)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "ok")
}

// LoginUser handles user login.
//
//	@Summary		Login User
//	@Description	Authenticate a user with email and password.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		params.LoginUserRequest	true	"Login User Request"
//	@Success		200		{string}	string					"ok"
//	@Failure		400		{object}	errs.ValidationError
//	@Failure		500		{object}	APIError
//	@Router			/auth/login [post]
func (ah *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	body := params.LoginUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, errs.ValidationError{Message: err.Error()})
		return
	}

	if err := body.Validate(); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	res, err := ah.svc.LoginUser(r.Context(), body)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	sendSuccessResponse(w, http.StatusOK, map[string]string{"token": res})
}
