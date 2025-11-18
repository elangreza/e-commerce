package client

import (
	"context"

	"github.com/elangreza/e-commerce/api/internal/params"
	"github.com/elangreza/e-commerce/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	productServiceClient struct {
		client gen.ProductServiceClient
	}
)

func NewProductClient() (*productServiceClient, error) {
	grpcClient, err := grpc.NewClient("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	productClient := gen.NewProductServiceClient(grpcClient)
	return &productServiceClient{client: productClient}, nil
}

func (s *productServiceClient) ListProducts(ctx context.Context, params params.ListProductsRequest) (*gen.ListProductsResponse, error) {
	req := &gen.ListProductsRequest{
		Search: "",
		Limit:  0,
		Page:   0,
		SortBy: "",
	}

	return s.client.ListProducts(ctx, req)
}
