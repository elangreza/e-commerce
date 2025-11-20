package handler

// go generate
//go:generate mockgen -source=product_grpc.go -destination=./mock/mock_product_grpc.go -package=mock

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/elangreza/e-commerce/api/internal/params"
	"github.com/elangreza/e-commerce/gen"
)

type (
	productService interface {
		ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error)
	}

	ProductHandler struct {
		productService productService
	}
)

func NewProductHandler(productService productService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
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

		products, err := s.productService.ListProducts(r.Context(), req)
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

func (s *ProductHandler) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.Product, error) {
	return nil, nil
}

func (s *ProductHandler) GetProducts(ctx context.Context, req *gen.GetProductsRequest) (*gen.Products, error) {
	return nil, nil
}
