package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"
	"github/elangreza/e-commerce/stock/internal/constanta"
	"github/elangreza/e-commerce/stock/internal/entity"
	"strings"
)

type StockRepo struct {
	db *sql.DB
}

func NewStockRepo(db *sql.DB) *StockRepo {
	return &StockRepo{db: db}
}

// GetStocks retrieves stock information for the given product IDs.
func (r *StockRepo) GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error) {
	if len(productIDs) == 0 {
		return []*entity.Stock{}, nil
	}

	placeholders := strings.Repeat("?,", len(productIDs))
	placeholders = strings.TrimRight(placeholders, ",")
	query := fmt.Sprintf(`SELECT product_id, quantity FROM stocks WHERE product_id IN (%s)`, placeholders)
	args := make([]any, len(productIDs))
	for i, id := range productIDs {
		args[i] = id
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []*entity.Stock
	for rows.Next() {
		var stock entity.Stock
		if err := rows.Scan(&stock.ProductID, &stock.Quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}

	return stocks, nil
}

func (r *StockRepo) ReserveStock(ctx context.Context, reserveStock entity.ReserveStock) ([]int64, error) {
	reservedStockIDs := []int64{}
	err := dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		for _, reqStock := range reserveStock.Stocks {
			var currQuantity int64
			err := tx.QueryRowContext(ctx, `
			SELECT SUM(quantity) as total_qty
				FROM stocks
			WHERE product_id = ?;`,
				reqStock.ProductID).Scan(&currQuantity)
			if err != nil {
				if err == sql.ErrNoRows {
					currQuantity = 0
				} else {
					return err
				}
			}

			if currQuantity == 0 {
				return fmt.Errorf("stock for product_id %s is empty", reqStock.ProductID)
			}

			if currQuantity < reqStock.Quantity {
				return fmt.Errorf("insufficient stock for product_id %s: requested %d, available %d", reqStock.ProductID, reqStock.Quantity, currQuantity)
			}

			currentStocks := []entity.Stock{}

			rows, err := tx.QueryContext(ctx, `
			SELECT 
				id, 
				quantity
			FROM (
				SELECT 
					id, 
					quantity, 
					created_at,
					SUM(quantity) OVER (ORDER BY created_at ASC) AS running_total
				FROM stocks
				WHERE product_id = ?
				ORDER BY created_at ASC
			) 
			WHERE running_total <= ? OR (
				running_total > ? AND (running_total - quantity) < ?
			)
			ORDER BY created_at ASC;
			`, reqStock.ProductID, reqStock.Quantity, reqStock.Quantity, reqStock.Quantity)

			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var stock entity.Stock
				if err := rows.Scan(&stock.ID, &stock.Quantity); err != nil {
					return err
				}
				currentStocks = append(currentStocks, stock)
			}
			if err := rows.Err(); err != nil {
				return err
			}

			// Allocate requested stock quantity by iterating through available stock entries (ordered by creation date).
			// For each stock entry, reserve as much as possible (up to the remaining requested quantity),
			// update the stock quantity, and record the reservation until the request is fulfilled.

			var currReqStock = reqStock.Quantity
			for _, currStock := range currentStocks {
				var qty = min(currStock.Quantity, currReqStock)

				_, err = tx.ExecContext(ctx, `UPDATE stocks SET quantity = quantity - ? WHERE id = ? AND quantity >= ?`, qty, currStock.ID, qty)
				if err != nil {
					return err
				}

				result, err := tx.ExecContext(ctx, `INSERT INTO reserved_stocks (stock_id, quantity, user_id, status) VALUES (?, ?, ?, ?)`, currStock.ID, qty, reserveStock.UserID, constanta.ReservedStockStatusReserved)
				if err != nil {
					return err
				}

				insertedID, err := result.LastInsertId()
				if err != nil {
					return err
				}

				reservedStockIDs = append(reservedStockIDs, insertedID)

				currReqStock -= qty
			}

		}
		return nil
	})

	if err != nil {
		return []int64{}, err
	}

	return reservedStockIDs, nil
}

func (r *StockRepo) ReleaseStock(ctx context.Context, releaseStock entity.ReleaseStock) ([]int64, error) {
	releasedStockIDs := []int64{}
	err := dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		for _, reservedStockID := range releaseStock.ReservedStockIDs {
			var quantity, stockID int
			err := tx.QueryRowContext(ctx, `SELECT quantity, stock_id FROM reserved_stocks WHERE id = ? AND user_id = ? AND status = ?`, reservedStockID, releaseStock.UserID, constanta.ReservedStockStatusReserved).Scan(&quantity, &stockID)
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, `UPDATE stocks SET quantity = quantity + ? WHERE id = ?`, quantity, stockID)
			if err != nil {
				return err
			}

			result, err := tx.ExecContext(ctx, `INSERT INTO released_stocks (stock_id, quantity, user_id, reserved_stock_id) VALUES (?, ?, ?, ?)`, stockID, quantity, releaseStock.UserID, reservedStockID)
			if err != nil {
				return err
			}

			insertedID, err := result.LastInsertId()
			if err != nil {
				return err
			}
			releasedStockIDs = append(releasedStockIDs, insertedID)

			_, err = tx.ExecContext(ctx, `UPDATE reserved_stocks SET status = ? WHERE id = ? AND status = ?`, constanta.ReservedStockStatusReleased, reservedStockID, constanta.ReservedStockStatusReserved)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return []int64{}, err
	}
	return releasedStockIDs, nil
}
