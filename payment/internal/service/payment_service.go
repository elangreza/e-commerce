package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/constanta"
	"github.com/elangreza/e-commerce/payment/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

type (
	paymentRepo interface {
		CreatePayment(ctx context.Context, payment entity.Payment) error
		UpdatePaymentStatusByTransactionID(ctx context.Context, paymentStatus constanta.PaymentStatus, transactionID string) error
		GetPaymentByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
		GetExpiryPayments(ctx context.Context, duration time.Duration) ([]entity.Payment, error)
	}
)

type paymentService struct {
	paymentRepo        paymentRepo
	maxTimeToBeExpired time.Duration
	orderService       gen.OrderServiceClient
	gen.UnimplementedPaymentServiceServer
}

func NewPaymentService(
	paymentRepo paymentRepo,
	maxTimeToBeExpired time.Duration,
	orderService gen.OrderServiceClient,
) *paymentService {
	return &paymentService{
		paymentRepo:        paymentRepo,
		maxTimeToBeExpired: maxTimeToBeExpired,
		orderService:       orderService,
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
	payment, err := p.paymentRepo.GetPaymentByTransactionID(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	if payment.Status != constanta.WAITING {
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

func (p *paymentService) GetPayment(ctx context.Context, req *gen.GetPaymentRequest) (*gen.GetPaymentResponse, error) {
	payment, err := p.paymentRepo.GetPaymentByTransactionID(ctx, req.TransactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "transaction not found")
		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &gen.GetPaymentResponse{
		TransactionId: payment.TransactionID,
		Status:        string(payment.Status),
		CreatedAt:     payment.CreatedAt.Format(time.DateTime),
		ExpiredAt:     payment.CreatedAt.Add(p.maxTimeToBeExpired).Format(time.DateTime),
		TotalAmount:   payment.TotalAmount,
	}, nil
}

func (p *paymentService) UpdatePayment(ctx context.Context, req *gen.UpdatePaymentRequest) (*gen.UpdatePaymentResponse, error) {
	payment, err := p.paymentRepo.GetPaymentByTransactionID(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	if payment.Status != constanta.WAITING {
		return &gen.UpdatePaymentResponse{
			Status: string(payment.Status),
		}, nil
	}

	if req.TotalAmount.Units > payment.TotalAmount.Units || req.TotalAmount.Units < payment.TotalAmount.Units {
		err = p.paymentRepo.UpdatePaymentStatusByTransactionID(ctx, constanta.FAILED, req.TransactionId)
		if err != nil {
			return nil, err
		}

		_, err = p.orderService.CallbackTransaction(ctx, &gen.CallbackTransactionRequest{
			TransactionId: req.TransactionId,
			PaymentStatus: constanta.FAILED.String(),
		})
		if err != nil {
			return nil, err
		}

		return &gen.UpdatePaymentResponse{
			Status: string(constanta.FAILED),
		}, nil
	}

	if req.TotalAmount.Units == payment.TotalAmount.Units {
		err = p.paymentRepo.UpdatePaymentStatusByTransactionID(ctx, constanta.PAID, req.TransactionId)
		if err != nil {
			return nil, err
		}

		_, err = p.orderService.CallbackTransaction(ctx, &gen.CallbackTransactionRequest{
			TransactionId: req.TransactionId,
			PaymentStatus: constanta.PAID.String(),
		})
		if err != nil {
			return nil, err
		}

		return &gen.UpdatePaymentResponse{
			Status: string(constanta.PAID),
		}, nil
	}

	return &gen.UpdatePaymentResponse{
		Status: "UNKNOWN",
	}, nil
}

func (p *paymentService) RemoveExpiryPayment(ctx context.Context, duration time.Duration) (int, error) {
	payments, err := p.paymentRepo.GetExpiryPayments(ctx, duration)
	if err != nil {
		return 0, err
	}

	for _, payment := range payments {
		if payment.Status != constanta.WAITING {
			continue
		}
		err = p.paymentRepo.UpdatePaymentStatusByTransactionID(ctx, constanta.FAILED, payment.TransactionID)
		if err != nil {
			fmt.Println("err when Update status", err)
		}

		_, err = p.orderService.CallbackTransaction(ctx, &gen.CallbackTransactionRequest{
			TransactionId: payment.TransactionID,
			PaymentStatus: constanta.FAILED.String(),
		})
		if err != nil {
			fmt.Println("err when callback transaction", err)
		}
	}

	return len(payments), nil
}
