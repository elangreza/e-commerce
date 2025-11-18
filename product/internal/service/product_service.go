package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/elangreza/e-commerce/product/internal/mockjson"
	params "github.com/elangreza/e-commerce/product/internal/params"
	"github.com/elangreza/e-commerce/product/pkg/errs"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	gen.UnimplementedProductServiceServer
}

func (p *productService) ListProducts(ctx context.Context, req *gen.ListProductsRequest) (*gen.ListProductsResponse, error) {
	paginationParams := params.PaginationParams{
		Sorts:  strings.Split(req.GetSortBy(), ","),
		Search: req.GetSearch(),
		Limit:  req.GetLimit(),
		Page:   req.GetPage(),
	}

	paginationParams.SetValidSortKey("updated_at", "name", "price")

	if err := paginationParams.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	reqParams := entity.ListProductRequest{
		Search:      paginationParams.Search,
		Page:        paginationParams.Page,
		Limit:       paginationParams.Limit,
		OrderClause: paginationParams.GetOrderClause(),
	}

	products, err := p.productRepo.ListProducts(ctx, reqParams)
	if err != nil {
		return nil, err
	}

	total, err := p.productRepo.TotalProducts(ctx, reqParams)
	if err != nil {
		return nil, err
	}

	productResponses := make([]*gen.Product, len(products))
	for i, product := range products {
		productResponses[i] = &gen.Product{
			Id:          product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
		}
	}

	return &gen.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		TotalPages: paginationParams.GetTotalPages(total),
	}, nil
}

func (p *productService) GetProduct(ctx context.Context, req *gen.GetProductRequest) (*gen.Product, error) {
	productID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	products, err := p.productRepo.GetProductByIDs(ctx, productID)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mockjson.DataNotFound) {
		// 	return nil, errs.NotFound{Message: "product not found"}
		// }
		return nil, err
	}

	if len(products) == 0 {
		return nil, errs.NotFound{Message: "product not found"}
	}

	stocks, err := p.stockServiceClient.GetStocks(ctx, []string{products[0].ID.String()})
	if err != nil {
		return nil, err
	}

	var stock int64 = 0
	for _, v := range stocks.Stocks {
		stock += v.Quantity
	}

	return &gen.Product{
		Id:          products[0].ID.String(), // Convert UUID to string
		Name:        products[0].Name,
		Description: products[0].Description,
		Price:       products[0].Price,
		ImageUrl:    products[0].ImageUrl,
		Stock:       stock,
	}, nil
}

func (p *productService) GetProducts(ctx context.Context, req *gen.GetProductsRequest) (*gen.Products, error) {
	productIDs := []uuid.UUID{}

	for _, productID := range req.Ids {
		pUUID, err := uuid.Parse(productID)
		if err != nil {
			return nil, err
		}
		productIDs = append(productIDs, pUUID)
	}

	products, err := p.productRepo.GetProductByIDs(ctx, productIDs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mockjson.DataNotFound) {
			return nil, errs.NotFound{Message: "product not found"}
		}
		return nil, err
	}

	res := []*gen.Product{}
	for _, product := range products {
		var stock int64 = 0
		if req.WithStock {
			stocks, err := p.stockServiceClient.GetStocks(ctx, []string{product.ID.String()})
			if err != nil {
				return nil, err
			}
			for _, v := range stocks.Stocks {
				stock += v.Quantity
			}
		}

		res = append(res, &gen.Product{
			Id:          product.ID.String(), // Convert UUID to string
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			Stock:       stock,
		})
	}

	return &gen.Products{
		Products: res,
	}, nil
}
