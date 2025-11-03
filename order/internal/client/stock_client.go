package client

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	stockServiceClient struct {
		client gen.StockServiceClient
	}
)

func NewStockClient() (*stockServiceClient, error) {
	grpcClient, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	stockClient := gen.NewStockServiceClient(grpcClient)
	return &stockServiceClient{client: stockClient}, nil
}

func (s *stockServiceClient) GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error) {
	return s.client.GetStocks(ctx, &gen.GetStockRequest{ProductIds: productIds})
}

// reserve stock after order is created
func (s *stockServiceClient) ReserveStock(ctx context.Context, cartItem []entity.CartItem) (*gen.ReserveStockResponse, error) {
	stocks := []*gen.Stock{}
	for _, item := range cartItem {
		stocks = append(stocks, &gen.Stock{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// add user id in context
	return s.client.ReserveStock(ctx, &gen.ReserveStockRequest{
		Stocks: stocks,
	})
}

// release stock when creating order is failed or order is cancelled
func (s *stockServiceClient) ReleaseStock(ctx context.Context, reservedStockIds []int64) (*gen.ReleaseStockResponse, error) {
	return s.client.ReleaseStock(ctx, &gen.ReleaseStockRequest{
		ReservedStockIds: reservedStockIds,
	})
}

// confirm stock after order is payed
func (s *stockServiceClient) ConfirmStock(ctx context.Context, productIds []string) (*gen.StockList, error) {
	return s.client.GetStocks(ctx, &gen.GetStockRequest{ProductIds: productIds})
}
