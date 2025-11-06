package sqlitedb

import (
	"context"
	"database/sql"

	"github/elangreza/e-commerce/pkg/dbsql"

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
		_, err := tx.ExecContext(ctx, `INSERT INTO orders(id, user_id, status, total_amount_units, currency) VALUES(?, ?, ?, ?, ?)`,
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
				item.PricePerUnit,
				item.Currency,
				item.Quantity,
				item.TotalPricePerUnit,
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
