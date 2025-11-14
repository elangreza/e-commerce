package sqlitedb

import (
	"context"
	"database/sql"
	"github/elangreza/e-commerce/pkg/money"
	"time"

	"github.com/elangreza/e-commerce/payment/internal/constanta"
	"github.com/elangreza/e-commerce/payment/internal/entity"
	"github.com/google/uuid"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (p *PaymentRepository) CreatePayment(ctx context.Context, payment entity.Payment) error {
	q := `INSERT INTO payments(id, status, total_amount, currency, transaction_id, order_id)
	VALUES (?, ?, ?, ?, ?, ?);`

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, q,
		id,
		payment.Status,
		payment.TotalAmount.Units,
		payment.TotalAmount.CurrencyCode,
		payment.TransactionID,
		payment.OrderID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PaymentRepository) UpdatePaymentStatusByTransactionID(ctx context.Context, paymentStatus constanta.PaymentStatus, transactionID string) error {
	q := `UPDATE payments
		SET status = ? AND updated_at = ?
		WHERE transaction_id = ?;`

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, q,
		id,
		paymentStatus,
		time.Now(),
		transactionID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PaymentRepository) GetPaymentStatusByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error) {
	q := `SELECT 
	id,
	status,
	total_amount,
	currency,
	transaction_id,
	order_id,
	created_at,
	updated_at
	FROM payments WHERE transaction_id = ?
	`

	var payment entity.Payment
	var totalAmount int64
	var currency string
	err := p.db.QueryRowContext(ctx, q, transactionID).Scan(
		&payment.ID,
		&payment.Status,
		&totalAmount,
		&currency,
		&payment.TransactionID,
		&payment.OrderID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	payment.TotalAmount, err = money.New(totalAmount, currency)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}
