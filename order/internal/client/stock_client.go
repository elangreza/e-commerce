package client

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
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
