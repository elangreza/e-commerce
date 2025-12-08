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
	OrderService interface {
		AddProductToCart(ctx context.Context, req params.AddToCartRequest) error
		GetCart(ctx context.Context) (*params.GetCartResponse, error)
		CreateOrder(ctx context.Context, req params.CreateOrderRequest) (*params.CreateOrderResponse, error)
	}

	orderHandler struct {
		svc OrderService
	}
)

func NewOrderHandler(
	publicRoute chi.Router,
	authService AuthService,
	svc OrderService,
) {

	authMiddleware := AuthMiddleware{
		svc: authService,
	}

	oh := orderHandler{
		svc: svc,
	}

	publicRoute.Group(func(r chi.Router) {
		r.Use(authMiddleware.MustAuthMiddleware())
		r.Post("/order", oh.CreateOrder())
		r.Post("/cart", oh.AddProductToCart())
		r.Get("/cart", oh.GetCart())
	})
}

func (oh *orderHandler) AddProductToCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := params.AddToCartRequest{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, errs.ValidationError{Message: err.Error()})
			return
		}

		if err := body.Validate(); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		ctx := r.Context()

		err := oh.svc.AddProductToCart(ctx, body)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		sendSuccessResponse(w, http.StatusCreated, "ok")
	}
}

func (oh *orderHandler) GetCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		cart, err := oh.svc.GetCart(ctx)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		sendSuccessResponse(w, http.StatusOK, cart)
	}
}

func (oh *orderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := params.CreateOrderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, errs.ValidationError{Message: err.Error()})
			return
		}

		if err := body.Validate(); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		ctx := r.Context()

		order, err := oh.svc.CreateOrder(ctx, body)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		sendSuccessResponse(w, http.StatusCreated, order)
	}
}
