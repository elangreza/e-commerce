package service

import (
	"context"

	"github.com/elangreza/e-commerce/product/internal/entity"
	params "github.com/elangreza/e-commerce/product/params"
	"github.com/google/uuid"
)

type (
	productRepo interface {
		ListProducts(ctx context.Context, req params.ListProductsRequest) ([]entity.Product, error)
		TotalProducts(ctx context.Context, req params.ListProductsRequest) (int64, error)
		GetProductByID(ctx context.Context, ID uuid.UUID) (*entity.Product, error)
	}

	ProductService interface {
		ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error)
		GetProduct(ctx context.Context, req params.GetProductRequest) (*params.GetProductResponse, error)
	}
)

func NewProductService(productRepo productRepo) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

type productService struct {
	productRepo productRepo
}

func (s *productService) ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error) {
	// Implementation for listing products
	products, err := s.productRepo.ListProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	productResponses := make([]params.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = params.ProductResponse{
			ID:          product.ID.String(), // Convert UUID to string
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Picture:     product.Picture,
		}
	}

	total, err := s.productRepo.TotalProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := total / req.Limit
	if total%req.Limit != 0 {
		totalPages++
	}

	return &params.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *productService) GetProduct(ctx context.Context, req params.GetProductRequest) (*params.GetProductResponse, error) {
	// Implementation for getting a product by ID
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &params.GetProductResponse{
		Product: &params.ProductResponse{
			ID:          product.ID.String(), // Convert UUID to string
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Picture:     product.Picture,
		},
	}, nil
}
