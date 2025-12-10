package client

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	paymentServiceClient struct {
		client gen.PaymentServiceClient
	}
)

func NewPaymentClient() (*paymentServiceClient, error) {
	grpcClient, err := grpc.NewClient("payment:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	stockClient := gen.NewPaymentServiceClient(grpcClient)
	return &paymentServiceClient{client: stockClient}, nil
}

func (s *paymentServiceClient) ProcessPayment(ctx context.Context, totalAmount *gen.Money, orderID uuid.UUID) (*gen.ProcessPaymentResponse, error) {
	return s.client.ProcessPayment(ctx, &gen.ProcessPaymentRequest{
		OrderId:     orderID.String(),
		TotalAmount: totalAmount,
	})
}
