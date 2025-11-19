package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/constanta"
	"github.com/elangreza/e-commerce/payment/internal/entity"
)

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

type (
	paymentRepo interface {
		CreatePayment(ctx context.Context, payment entity.Payment) error
		UpdatePaymentStatusByTransactionID(ctx context.Context, paymentStatus constanta.PaymentStatus, transactionID string) error
		GetPaymentStatusByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
	}
)

type paymentService struct {
	paymentRepo paymentRepo
	gen.UnimplementedPaymentServiceServer
}

func NewPaymentService(
	paymentRepo paymentRepo,
) *paymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
	}
}

func (p *paymentService) ProcessPayment(ctx context.Context, req *gen.ProcessPaymentRequest) (*gen.ProcessPaymentResponse, error) {
	transactionID := generateBase62ID(defaultLength)
	err := p.paymentRepo.CreatePayment(ctx, entity.Payment{
		Status:        constanta.WAITING,
		TotalAmount:   req.TotalAmount,
		TransactionID: transactionID,
		OrderID:       req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	return &gen.ProcessPaymentResponse{
		TransactionId: transactionID,
	}, nil
}

func (p *paymentService) RollbackPayment(ctx context.Context, req *gen.RollbackPaymentRequest) (*gen.Empty, error) {
	payment, err := p.paymentRepo.GetPaymentStatusByTransactionID(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	if !(payment.Status == constanta.WAITING) {
		return nil, fmt.Errorf("payment must be waiting rollback the payment")
	}

	err = p.paymentRepo.UpdatePaymentStatusByTransactionID(ctx, constanta.CANCELLED, req.TransactionId)
	if err != nil {
		return nil, err
	}

	return &gen.Empty{}, nil
}

const (
	defaultLength = 8
	base62Chars   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func generateBase62ID(length int) string {
	result := make([]byte, length)
	for i := range length {
		num, _ := rand.Int(rand.Reader, big.NewInt(62))
		result[i] = base62Chars[num.Int64()]
	}
	return string(result)
}
