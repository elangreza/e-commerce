package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"
	"database/sql"
	"errors"

	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/elangreza/e-commerce/product/internal/mockjson"
	params "github.com/elangreza/e-commerce/product/params"
	"github.com/elangreza/e-commerce/product/pkg/errs"
	"github.com/google/uuid"
)

type (
	productRepo interface {
		ListProducts(ctx context.Context, req entity.ListProductRequest) ([]entity.Product, error)
		TotalProducts(ctx context.Context, req entity.ListProductRequest) (int64, error)
		GetProductByID(ctx context.Context, ID uuid.UUID) (*entity.Product, error)
	}
)

func NewProductService(productRepo productRepo) *productService {
	return &productService{
		productRepo: productRepo,
	}
}

type productService struct {
	productRepo productRepo
}

func (s *productService) ListProducts(ctx context.Context, req params.PaginationRequest) (*params.ListProductsResponse, error) {
	// Implementation for listing products
	reqParams := entity.ListProductRequest{
		Search: req.Search,
		Page:   req.Page,
		Limit:  req.Limit,
		SortBy: req.SortBy,
	}
	products, err := s.productRepo.ListProducts(ctx, reqParams)
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

	total, err := s.productRepo.TotalProducts(ctx, reqParams)
	if err != nil {
		return nil, err
	}

	totalPages := total / reqParams.Limit
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
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mockjson.DataNotFound) {
			return nil, errs.NotFound{Message: "product not found"}
		}
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
