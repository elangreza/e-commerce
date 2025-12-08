package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elangreza/e-commerce/api/internal/params"
	"github.com/go-chi/chi/v5"
)

type (
	ProductService interface {
		ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error)
	}

	ProductHandler struct {
		svc ProductService
	}
)

func NewProductHandler(ar chi.Router, ps ProductService) {

	authHandler := ProductHandler{
		svc: ps,
	}

	ar.Route("/products", func(r chi.Router) {
		r.Get("/", authHandler.ListProducts())
	})
}

func (s *ProductHandler) ListProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req params.ListProductsRequest

		queries := r.URL.Query()

		req.Search = queries.Get("search")
		req.SortBy = queries.Get("sort_by")
		if len(queries["limit"]) > 0 {
			limit, _ := strconv.Atoi(queries["limit"][0])
			req.Limit = int64(limit)
		}

		if len(queries["page"]) > 0 {
			page, _ := strconv.Atoi(queries["page"][0])
			req.Page = int64(page)
		}

		products, err := s.svc.ListProducts(r.Context(), req)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		sendSuccessResponse(w, http.StatusOK, products)
	}
}
