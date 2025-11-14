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
}

func NewPaymentService(
	paymentRepo paymentRepo,
) *paymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
	}
}

func (p *paymentService) ProcessPayment(ctx context.Context, totalAmount *gen.Money, orderID string) (string, error) {
	transactionID := generateBase62ID(defaultLength)
	err := p.paymentRepo.CreatePayment(ctx, entity.Payment{
		Status:        constanta.WAITING,
		TotalAmount:   totalAmount,
		TransactionID: transactionID,
		OrderID:       orderID,
	})
	if err != nil {
		return "", err
	}

	return transactionID, nil
}

func (p *paymentService) RollbackPayment(ctx context.Context, transactionID string) error {
	payment, err := p.paymentRepo.GetPaymentStatusByTransactionID(ctx, transactionID)
	if err != nil {
		return err
	}

	if !(payment.Status == constanta.WAITING) {
		return fmt.Errorf("payment must be waiting rollback the payment")
	}

	err = p.paymentRepo.UpdatePaymentStatusByTransactionID(ctx, constanta.CANCELLED, transactionID)
	if err != nil {
		return err
	}

	return nil
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
