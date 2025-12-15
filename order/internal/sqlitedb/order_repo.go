package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/money"

	"github.com/elangreza/e-commerce/order/internal/constanta"
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
			) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
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

		_, err = tx.ExecContext(ctx, "UPDATE carts SET is_active = FALSE WHERE user_id = ?", order.UserID)
		if err != nil {
			return err
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
	transaction_id,
	created_at, 
	updated_at FROM orders WHERE idempotency_key = ?;`

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
		&ord.TransactionID,
		&ord.CreatedAt,
		&ord.UpdatedAt,
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

func (r *OrderRepository) UpdateOrder(ctx context.Context, payloads map[string]any, orderID uuid.UUID) error {

	q := `UPDATE orders SET `
	fields := map[string]struct{}{
		"transaction_id": {},
		"status":         {},
	}
	args := []any{}
	for key, payload := range payloads {
		if _, ok := fields[key]; ok {
			q += fmt.Sprintf(" %s = ?,", key)
		}
		args = append(args, payload)
	}
	q += "updated_at = ? WHERE id = ?;"
	args = append(args, time.Now(), orderID)

	_, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetExpiryOrders(ctx context.Context, duration time.Duration) ([]entity.Order, error) {
	q := `SELECT 
	id,
	status,
	user_id
	FROM orders WHERE (status = ? OR status = ?) AND created_at < DATETIME(?);`

	timeLimit := time.Now().UTC().Add(-duration)

	rows, err := r.db.QueryContext(ctx,
		q,
		constanta.OrderStatusPending,
		constanta.OrderStatusStockReserved,
		timeLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []entity.Order{}
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.Status, &order.UserID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateOrderStatusWithCallback(ctx context.Context, status constanta.OrderStatus, orderID uuid.UUID, callback func() error) error {
	err := dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE orders SET status = ? WHERE id = ?`, status, orderID)
		if err != nil {
			return err
		}

		err = callback()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetOrderByTransactionID(ctx context.Context, transactionID string) (*entity.Order, error) {
	q := `SELECT id, status FROM orders WHERE transaction_id = ?;`

	var ord entity.Order
	err := r.db.QueryRowContext(ctx, q, transactionID).Scan(
		&ord.ID,
		&ord.Status,
	)
	if err != nil {
		return nil, err
	}

	return &ord, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.Order, error) {
	q := `SELECT id, 
	idempotency_key, 
	user_id, 
	status, 
	total_amount, 
	currency, 
	transaction_id, 
	created_at,
	updated_at FROM orders WHERE id = ?;`

	var totalAmount int64
	var currencyCode string
	var ord entity.Order
	err := r.db.QueryRowContext(ctx, q, orderID).Scan(
		&ord.IdempotencyKey,
		&ord.ID,
		&ord.UserID,
		&ord.Status,
		&totalAmount,
		&currencyCode,
		&ord.TransactionID,
		&ord.CreatedAt,
		&ord.UpdatedAt,
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

func (r *OrderRepository) GetOrderList(ctx context.Context, req entity.GetOrderListRequest) ([]entity.Order, error) {
	q := `SELECT id, 
	idempotency_key, 
	user_id, 
	status, 
	total_amount, 
	currency, 
	transaction_id,
	created_at, 
	updated_at FROM orders WHERE user_id = ?;`

	if req.IsFilterByDate {
		q += " AND created_at BETWEEN ? AND ?;"
	}

	if req.IsFilterByStatus {
		q += " AND status = ?;"
	}

	args := []any{req.UserID}
	if req.IsFilterByDate {
		args = append(args, req.StartDate, req.EndDate)
	}
	if req.IsFilterByStatus {
		args = append(args, req.Status)
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []entity.Order{}
	for rows.Next() {
		var order entity.Order
		var totalAmount int64
		var currencyCode string
		err := rows.Scan(
			&order.IdempotencyKey,
			&order.ID,
			&order.UserID,
			&order.Status,
			&totalAmount,
			&currencyCode,
			&order.TransactionID,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		order.TotalAmount, err = money.New(totalAmount, currencyCode)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
