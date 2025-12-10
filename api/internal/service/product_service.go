package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"

	params "github.com/elangreza/e-commerce/api/internal/params"
	"github.com/elangreza/e-commerce/gen"
)

func NewProductService(
	pClient gen.ProductServiceClient,
	sClient gen.ShopServiceClient,
) *productService {
	return &productService{
		productServiceClient: pClient,
		shopServiceClient:    sClient,
	}
}

type productService struct {
	productServiceClient gen.ProductServiceClient
	shopServiceClient    gen.ShopServiceClient
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

	shopIDs := []int64{}
	for _, product := range listProduct.Products {
		shopIDs = append(shopIDs, product.ShopId)
	}

	shopMap, err := s.getShopMap(ctx, shopIDs)
	if err != nil {
		return nil, convertErrGrpc(err)
	}

	for _, product := range listProduct.Products {
		p := &params.Product{
			Id:          product.GetId(),
			Name:        product.GetName(),
			Description: product.GetDescription(),
			ImageUrl:    product.GetImageUrl(),
			Stock:       product.GetStock(),
			Price: &params.Money{
				Units:        product.Price.GetUnits(),
				CurrencyCode: product.Price.GetCurrencyCode(),
			},
			ShopID: product.GetShopId(),
		}
		shopName, ok := shopMap[product.GetShopId()]
		if ok {
			p.ShopName = shopName
		}
		res.Products = append(res.Products, p)
	}

	return res, nil
}

func (s *productService) GetProductsDetails(ctx context.Context, req params.GetProductsDetail) (*params.ListProductsResponse, error) {
	listProduct, err := s.productServiceClient.GetProducts(ctx, &gen.GetProductsRequest{
		Ids:       req.Ids,
		WithStock: req.WithStock,
	})
	if err != nil {
		return nil, convertErrGrpc(err)
	}

	shopIDs := []int64{}
	for _, product := range listProduct.Products {
		shopIDs = append(shopIDs, product.ShopId)
	}

	shopMap, err := s.getShopMap(ctx, shopIDs)
	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.ListProductsResponse{
		Products: []*params.Product{},
	}

	for _, product := range listProduct.Products {
		p := &params.Product{
			Id:          product.GetId(),
			Name:        product.GetName(),
			Description: product.GetDescription(),
			ImageUrl:    product.GetImageUrl(),
			Price: &params.Money{
				Units:        product.Price.GetUnits(),
				CurrencyCode: product.Price.GetCurrencyCode(),
			},
			ShopID: product.ShopId,
		}
		if req.WithStock {
			p.Stock = product.GetStock()
		}

		shopName, ok := shopMap[product.GetShopId()]
		if ok {
			p.ShopName = shopName
		}
		res.Products = append(res.Products, p)
	}

	return res, nil
}

func (s *productService) getShopMap(ctx context.Context, shopIDs []int64) (map[int64]string, error) {
	shops, err := s.shopServiceClient.GetShops(ctx, &gen.GetShopsRequest{
		Ids:            shopIDs,
		WithWarehouses: false,
	})
	if err != nil {
		return nil, err
	}

	shopsMap := make(map[int64]string)
	for _, shop := range shops.Shops {
		shopsMap[shop.GetId()] = shop.Name
	}

	return shopsMap, nil
}
