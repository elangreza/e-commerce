package client

import (
	"context"

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
	grpcClient, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	productClient := gen.NewProductServiceClient(grpcClient)
	return &productServiceClient{client: productClient}, nil
}

func (s *productServiceClient) GetProduct(ctx context.Context, productId string) (*gen.Product, error) {
	return s.client.GetProduct(ctx, &gen.GetProductRequest{Id: productId})
}

func (s *productServiceClient) GetProducts(ctx context.Context, withStock bool, productIds ...string) (*gen.Products, error) {
	return s.client.GetProducts(ctx, &gen.GetProductsRequest{
		Ids:       productIds,
		WithStock: withStock,
	})
}
