package service

//go:generate mockgen -source=product_service.go -destination=mock/mock_product_service.go -package=mock
//go:generate mockgen -package=mock -destination=mock/mock_deps.go github.com/elangreza/e-commerce/gen WarehouseServiceClient

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/entity"
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
)

func NewProductService(productRepo productRepo, warehouseServiceClient gen.WarehouseServiceClient) *ProductService {
	return &ProductService{
		productRepo:            productRepo,
		warehouseServiceClient: warehouseServiceClient,
	}
}

type ProductService struct {
	productRepo            productRepo
	warehouseServiceClient gen.WarehouseServiceClient
	gen.UnimplementedProductServiceServer
}

func (p *ProductService) ListProducts(ctx context.Context, req *gen.ListProductsRequest) (*gen.ListProductsResponse, error) {
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

	var stockMap map[string]int64
	if req.WithStock {
		stockMap, err = p.getStockMap(ctx, products)
		if err != nil {
			return nil, err
		}
	}

	productResponses := make([]*gen.Product, len(products))
	for i, product := range products {
		var stock int64 = 0
		if req.WithStock {
			stk, ok := stockMap[product.ID.String()]
			if ok {
				stock = stk
			}
		}

		productResponses[i] = &gen.Product{
			Id:          product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			Stock:       stock,
			ShopId:      product.ShopID,
		}
	}

	return &gen.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		TotalPages: paginationParams.GetTotalPages(total),
	}, nil
}

func (p *ProductService) GetProducts(ctx context.Context, req *gen.GetProductsRequest) (*gen.Products, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{Message: "product not found"}
		}
		return nil, err
	}

	var stockMap map[string]int64
	if req.WithStock {
		stockMap, err = p.getStockMap(ctx, products)
		if err != nil {
			return nil, err
		}
	}

	productResponses := make([]*gen.Product, len(products))
	for i, product := range products {
		var stock int64 = 0
		if req.WithStock {
			stk, ok := stockMap[product.ID.String()]
			if ok {
				stock = stk
			}
		}

		productResponses[i] = &gen.Product{
			Id:          product.ID.String(), // Convert UUID to string
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageUrl:    product.ImageUrl,
			Stock:       stock,
			ShopId:      product.ShopID,
		}
	}

	return &gen.Products{
		Products: productResponses,
	}, nil
}

func (p *ProductService) getStockMap(ctx context.Context, products []entity.Product) (map[string]int64, error) {
	if len(products) == 0 {
		return nil, nil
	}

	productIDs := []string{}
	for _, p := range products {
		productIDs = append(productIDs, p.ID.String())
	}

	stocks, err := p.warehouseServiceClient.GetStocks(ctx,
		&gen.GetStockRequest{
			ProductIds: productIDs,
		})
	if err != nil {
		return nil, err
	}

	res := make(map[string]int64)
	for _, v := range stocks.Stocks {
		var stock int64
		stock += v.Quantity
		val, ok := res[v.ProductId]
		if ok {
			res[v.ProductId] = val + stock
		} else {
			res[v.ProductId] = stock
		}
	}

	return res, nil
}
