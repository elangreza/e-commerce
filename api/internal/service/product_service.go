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

func NewProductService(pClient gen.ProductServiceClient) *productService {
	return &productService{
		productServiceClient: pClient,
	}
}

type productService struct {
	productServiceClient gen.ProductServiceClient
}

func (s *productService) ListProducts(ctx context.Context, req params.ListProductsRequest) (*params.ListProductsResponse, error) {

	listProduct, err := s.productServiceClient.ListProducts(ctx, &gen.ListProductsRequest{
		Search: req.Search,
		Limit:  req.Limit,
		Page:   req.Page,
		SortBy: req.SortBy,
	})

	if err != nil {
		return nil, err
	}

	res := &params.ListProductsResponse{
		Products:   []*params.Product{},
		Total:      listProduct.GetTotal(),
		TotalPages: listProduct.GetTotalPages(),
	}

	for _, product := range listProduct.Products {
		res.Products = append(res.Products, &params.Product{
			Id:          product.GetId(),
			Name:        product.GetName(),
			Description: product.GetDescription(),
			ImageUrl:    product.GetImageUrl(),
			Price: &params.Money{
				Units:        product.Price.GetUnits(),
				CurrencyCode: product.Price.GetCurrencyCode(),
			},
		})
	}

	return res, nil
}
