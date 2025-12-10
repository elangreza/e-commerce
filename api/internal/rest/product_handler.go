package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/elangreza/e-commerce/api/internal/params"
	"github.com/go-chi/chi/v5"
)

type (
	ProductService interface {
		ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error)
		GetProductsDetails(ctx context.Context, req params.GetProductsDetail) (*params.ListProductsResponse, error)
	}

	ProductHandler struct {
		svc ProductService
	}
)

func NewProductHandler(ar chi.Router, ps ProductService) {

	authHandler := ProductHandler{
		svc: ps,
	}

	ar.Get("/products", authHandler.ListProducts())
	ar.Get("/product", authHandler.GetProductsDetails())

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

func (s *ProductHandler) GetProductsDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req params.GetProductsDetail

		queries := r.URL.Query()
		ids, ok := queries["id"]
		if !ok {
			err := errors.New("must provide id in get products")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Ids = append(req.Ids, ids...)

		var err error
		var withStock bool

		if queries.Has("with_stock") {
			withStockParams := queries.Get("with_stock")
			withStock, err = strconv.ParseBool(withStockParams)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		req.WithStock = withStock

		products, err := s.svc.GetProductsDetails(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
