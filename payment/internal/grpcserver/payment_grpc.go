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
	paymentService interface {
		ProcessPayment(ctx context.Context, totalAmount *gen.Money, orderID string) (string, error)
		RollbackPayment(ctx context.Context, transactionID string) error
	}

	PaymentServer struct {
		paymentService paymentService
		gen.UnimplementedPaymentServiceServer
	}
)

func NewPaymentServer(ps paymentService) *PaymentServer {
	return &PaymentServer{
		paymentService: ps,
	}
}

func (p *PaymentServer) ProcessPayment(ctx context.Context, req *gen.ProcessPaymentRequest) (*gen.ProcessPaymentResponse, error) {
	transactionID, err := p.paymentService.ProcessPayment(ctx, req.TotalAmount, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.ProcessPaymentResponse{
		TransactionId: transactionID,
	}, nil
}

func (p *PaymentServer) RollbackPayment(ctx context.Context, req *gen.RollbackPaymentRequest) (*gen.Empty, error) {

	err := p.paymentService.RollbackPayment(ctx, req.TransactionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.Empty{}, nil
}
