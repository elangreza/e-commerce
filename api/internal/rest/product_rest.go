package grpcserver

// go generate
//go:generate mockgen -source=product_grpc.go -destination=./mock/mock_product_grpc.go -package=mock

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
)

type (
	productService interface {
		// ListProducts(ctx context.Context, req params.PaginationParams) (*params.ListProductsResponse, error)
		// GetProduct(ctx context.Context, req params.GetProductRequest) (*params.GetProductResponse, error)
		// GetProducts(ctx context.Context, req params.GetProductsRequest) (*params.GetProductsResponse, error)
	}

	ProductServer struct {
		productService productService
		gen.UnimplementedProductServiceServer
	}
)

func NewProductServer(productService productService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

func (s *ProductServer) ListProducts(ctx context.Context, req *gen.ListProductsRequest) (*gen.ListProductsResponse, error) {
	return nil, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.Product, error) {
	return nil, nil
}

func (s *ProductServer) GetProducts(ctx context.Context, req *gen.GetProductsRequest) (*gen.Products, error) {
	return nil, nil
}
