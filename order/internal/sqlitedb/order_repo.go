package sqlitedb

import (
	"context"
	"database/sql"

	"github/elangreza/e-commerce/pkg/dbsql"
	"github/elangreza/e-commerce/pkg/money"

	"github.com/elangreza/e-commerce/order/internal/entity"
	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order entity.Order) (uuid.UUID, error) {
	// Implementation to create a new Order in the database

	orderID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	err = dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO orders(idempotency_key, id, user_id, status, total_amount, currency) VALUES(?, ?, ?, ?, ?, ?)`,
			order.IdempotencyKey,
			orderID,
			order.UserID,
			order.Status,
			order.TotalAmount.Units,
			order.TotalAmount.CurrencyCode,
		)
		if err != nil {
			return err
		}

		for _, item := range order.Items {

			orderItemID, err := uuid.NewV7()
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, `INSERT INTO order_items(
			    id,
				order_id,
				product_id,
				name,
				price_per_unit_units,
				currency,
				quantity,
				total_price_units
			) VALUES(?, ?, ?, ?, ?)`,
				orderItemID,
				orderID,
				item.ProductID,
				item.Name,
				item.PricePerUnit.GetUnits(),
				item.PricePerUnit.GetCurrencyCode(),
				item.Quantity,
				item.TotalPricePerUnit.GetUnits(),
			)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

func (r *OrderRepository) GetOrderByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*entity.Order, error) {
	q := `SELECT id, 
	idempotency_key, 
	user_id, 
	status, 
	total_amount, 
	currency, 
	created_at, 
	updated_at, 
	shipped_at, 
	cancelled_at FROM orders WHERE idempotency_key = ?;`

	var totalAmount int64
	var currencyCode string
	var ord entity.Order
	err := r.db.QueryRowContext(ctx, q, idempotencyKey).Scan(
		&ord.IdempotencyKey,
		&ord.ID,
		&ord.UserID,
		&ord.Status,
		&totalAmount,
		&currencyCode,
		&ord.CreatedAt,
		&ord.UpdatedAt,
		&ord.ShippedAt,
		&ord.CancelledAt,
	)
	if err != nil {
		return nil, err
	}

	ord.TotalAmount, err = money.New(totalAmount, currencyCode)
	if err != nil {
		return nil, err
	}

	qItems := `SELECT 
	id, 
	order_id, 
	product_id, 
	name, 
	price_per_unit_units, 
	currency, 
	quantity, 
	total_price_units
	FROM order_items WHERE order_id = ?;`

	rows, err := r.db.QueryContext(ctx, qItems, ord.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderItem entity.OrderItem
		var pricePerUnit int64
		var totalPricePerUnit int64
		var currencyCode string
		err = rows.Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.Name,
			&pricePerUnit,
			&currencyCode,
			&orderItem.Quantity,
			&totalPricePerUnit,
		)
		if err != nil {
			return nil, err
		}

		orderItem.PricePerUnit, err = money.New(pricePerUnit, currencyCode)
		if err != nil {
			return nil, err
		}

		orderItem.TotalPricePerUnit, err = money.New(totalPricePerUnit, currencyCode)
		if err != nil {
			return nil, err
		}

		ord.Items = append(ord.Items, orderItem)
	}

	return &ord, nil
}
