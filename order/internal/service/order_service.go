package service

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
)

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

type (
	orderRepo interface {
	}

	stockServiceClient interface {
		GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error)
	}
)

func NewOrderService(orderRepo orderRepo, stockServiceClient stockServiceClient) *orderService {
	return &orderService{
		orderRepo:          orderRepo,
		stockServiceClient: stockServiceClient,
	}
}

type orderService struct {
	orderRepo          orderRepo
	stockServiceClient stockServiceClient
}
