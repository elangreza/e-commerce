package grpcserver

// go generate
//go:generate mockgen -source=order_grpc.go -destination=./mock/mock_order_grpc.go -package=mock

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	orderService interface {
		ProcessPayment(ctx context.Context, totalAmount *gen.Money, orderID string) (string, error)
	}

	OrderServer struct {
		orderService orderService
		gen.UnimplementedPaymentServiceServer
	}
)

func NewPaymentServer(orderService orderService) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}

func (o *OrderServer) ProcessPayment(ctx context.Context, req *gen.ProcessPaymentRequest) (*gen.ProcessPaymentResponse, error) {
	transactionID, err := o.orderService.ProcessPayment(ctx, req.TotalAmount, req.OrderId)
	if err != nil {
		return &gen.ProcessPaymentResponse{
			TransactionId: "",
			Success:       false,
			ErrorMessage:  err.Error(),
		}, nil
	}

	return &gen.ProcessPaymentResponse{
		TransactionId: transactionID,
		Success:       true,
		ErrorMessage:  "",
	}, nil
}

func (o *OrderServer) RollbackPayment(ctx context.Context, req *gen.RollbackPaymentRequest) (*gen.RollbackPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RollbackPayment not implemented")
}
