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

func NewWarehouseClient(addr string) (*warehouseServiceClient, error) {
	grpcClient, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	warehouseClient := gen.NewWarehouseServiceClient(grpcClient)
	return &warehouseServiceClient{client: warehouseClient}, nil
}

// GetStocks implements StockServiceClient.
func (s *warehouseServiceClient) GetWarehouseByShopID(ctx context.Context, shopID int64) (*gen.GetWarehouseByShopIDResponse, error) {
	return s.client.GetWarehouseByShopID(ctx, &gen.GetWarehouseByShopIDRequest{
		ShopId: shopID,
	})
}
