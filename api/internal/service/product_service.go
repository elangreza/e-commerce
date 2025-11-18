package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"

	params "github.com/elangreza/e-commerce/api/internal/params"
	"github.com/elangreza/e-commerce/gen"
)

type (
	productServiceClient interface {
		ListProducts(ctx context.Context, params params.ListProductsRequest) (*gen.ListProductsResponse, error)
	}
)

func NewProductService(stockServiceClient productServiceClient) *productService {
	return &productService{
		productServiceClient: stockServiceClient,
	}
}

type productService struct {
	productServiceClient productServiceClient
}

func (s *productService) ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error) {
	// products, err := s.productServiceClient.ListProducts(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
