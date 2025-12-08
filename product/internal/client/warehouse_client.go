package client

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	warehouseServiceClient struct {
		client gen.WarehouseServiceClient
	}
)

func NewWarehouseClient() (*warehouseServiceClient, error) {
	grpcClient, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	warehouseClient := gen.NewWarehouseServiceClient(grpcClient)
	return &warehouseServiceClient{client: warehouseClient}, nil
}

// GetStocks implements StockServiceClient.
func (s *warehouseServiceClient) GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error) {
	return s.client.GetStocks(ctx, &gen.GetStockRequest{ProductIds: productIds})
}
