package client

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/entity"
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

// GetStocks implements StockServiceClient.
func (s *stockServiceClient) GetStocks(ctx context.Context, productIds []string) ([]entity.Stock, error) {
	resp, err := s.client.GetStocks(ctx, &gen.GetStockRequest{ProductIds: productIds})
	if err != nil {
		return nil, err
	}

	var stocks []entity.Stock
	for _, item := range resp.Stocks {
		stocks = append(stocks, entity.Stock{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return stocks, nil
}
