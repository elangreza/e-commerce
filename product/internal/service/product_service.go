package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"
	"database/sql"
	"errors"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/elangreza/e-commerce/product/internal/mockjson"
	params "github.com/elangreza/e-commerce/product/internal/params"
	"github.com/elangreza/e-commerce/product/pkg/errs"
	"github.com/google/uuid"
)

type (
	productRepo interface {
		ListProducts(ctx context.Context, req entity.ListProductRequest) ([]entity.Product, error)
		TotalProducts(ctx context.Context, req entity.ListProductRequest) (int64, error)
		GetProductByIDs(ctx context.Context, ID ...uuid.UUID) ([]entity.Product, error)
	}

	stockServiceClient interface {
		GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error)
	}
)

func NewProductService(productRepo productRepo, stockServiceClient stockServiceClient) *productService {
	return &productService{
		productRepo:        productRepo,
		stockServiceClient: stockServiceClient,
	}
}

type productService struct {
	productRepo        productRepo
	stockServiceClient stockServiceClient
}

func (s *productService) ListProducts(ctx context.Context, req params.PaginationParams) (*params.ListProductsResponse, error) {
	reqParams := entity.ListProductRequest{
		Search:      req.Search,
		Page:        req.Page,
		Limit:       req.Limit,
		OrderClause: req.GetOrderClause(),
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
			ImageUrl:    product.ImageUrl,
		}
	}

	total, err := s.productRepo.TotalProducts(ctx, reqParams)
	if err != nil {
		return nil, err
	}

	return &params.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		TotalPages: req.GetTotalPages(total),
	}, nil
}

func (s *productService) GetProduct(ctx context.Context, req params.GetProductRequest) (*params.GetProductResponse, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, err
	}

	products, err := s.productRepo.GetProductByIDs(ctx, productID)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mockjson.DataNotFound) {
		// 	return nil, errs.NotFound{Message: "product not found"}
		// }
		return nil, err
	}

	if len(products) == 0 {
		return nil, errs.NotFound{Message: "product not found"}
	}

	stocks, err := s.stockServiceClient.GetStocks(ctx, []string{products[0].ID.String()})
	if err != nil {
		return nil, err
	}

	var stock int64 = 0
	for _, v := range stocks.Stocks {
		stock += v.Quantity
	}

	return &params.GetProductResponse{
		Product: &params.ProductResponse{
			ID:          products[0].ID.String(), // Convert UUID to string
			Name:        products[0].Name,
			Description: products[0].Description,
			Price:       products[0].Price,
			ImageUrl:    products[0].ImageUrl,
			Stock:       stock,
		},
	}, nil
}

func (s *productService) GetProducts(ctx context.Context, req params.GetProductsRequest) (*params.GetProductsResponse, error) {
	productIDs := []uuid.UUID{}

	for _, productID := range req.ProductIDs {
		pUUID, err := uuid.Parse(productID)
		if err != nil {
			return nil, err
		}
		productIDs = append(productIDs, pUUID)
	}

	products, err := s.productRepo.GetProductByIDs(ctx, productIDs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mockjson.DataNotFound) {
			return nil, errs.NotFound{Message: "product not found"}
		}
		return nil, err
	}

	res := []params.ProductResponse{}
	for _, product := range products {
		var stock int64 = 0
		if req.WithStock {
			stocks, err := s.stockServiceClient.GetStocks(ctx, []string{product.ID.String()})
			if err != nil {
				return nil, err
			}
			for _, v := range stocks.Stocks {
				stock += v.Quantity
			}
		}

		res = append(res, params.ProductResponse{
			ID:          product.ID.String(), // Convert UUID to string
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			Stock:       stock,
		})
	}

	return &params.GetProductsResponse{
		Products: res,
	}, nil
}
