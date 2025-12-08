package sqlitedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/elangreza/e-commerce/warehouse/internal/constanta"
	"github.com/elangreza/e-commerce/warehouse/internal/entity"

	"github.com/elangreza/e-commerce/pkg/dbsql"

	"github.com/google/uuid"
)

type WarehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepo(db *sql.DB) *WarehouseRepo {
	return &WarehouseRepo{db: db}
}

// GetStocks retrieves stock information for the given product IDs.
func (r *WarehouseRepo) GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error) {
	if len(productIDs) == 0 {
		return []*entity.Stock{}, nil
	}

	placeholders := strings.Repeat("?,", len(productIDs))
	placeholders = strings.TrimRight(placeholders, ",")
	query := fmt.Sprintf(`SELECT 
		s.product_id, s.quantity 
		FROM 
			stocks s
		LEFT JOIN warehouses w ON w.id = s.warehouse_id 
		WHERE 
			w.is_active IS TRUE AND s.product_id IN (%s) AND s.quantity > 0`, placeholders)
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

func (r *WarehouseRepo) ReserveStock(ctx context.Context, reserveStock entity.ReserveStock) ([]int64, error) {
	reservedStockIDs := []int64{}
	err := dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		for _, reqStock := range reserveStock.Stocks {
			var currQuantity int64
			err := tx.QueryRowContext(ctx, `
			SELECT SUM(s.quantity) as total_qty
				FROM stocks s
			LEFT JOIN warehouses w ON w.id = s.warehouse_id 
			WHERE w.is_active IS TRUE AND s.product_id = ?;`,
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
						s.id, 
						s.quantity, 
						s.created_at,
						SUM(s.quantity) OVER (ORDER BY s.created_at, s.id ASC) AS running_total
					FROM stocks s
					LEFT JOIN warehouses w ON w.id = s.warehouse_id 
					WHERE w.is_active IS TRUE AND s.product_id = ?
				) s
				WHERE 
					running_total <= ? 
				OR (
					running_total > ? 
					AND (running_total - quantity) < ?
				)
				ORDER BY created_at, id;
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

			// this reserve stock is using FIFO method
			// meaning the oldest stock will be reserved first
			// this is done to prevent stock from being expired
			// before it is sold
			// so we need to reserve the oldest stock first.
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

				result, err := tx.ExecContext(ctx, `INSERT INTO reserved_stocks (stock_id, quantity, user_id, status, order_id) VALUES (?, ?, ?, ?, ?)`,
					currStock.ID,
					qty,
					reserveStock.UserID,
					constanta.ReservedStockStatusReserved,
					reserveStock.OrderID)
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

func (r *WarehouseRepo) ReleaseStock(ctx context.Context, releaseStock entity.ReleaseStock) ([]int64, error) {
	releasedStockIDs := []int64{}
	err := dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {

		reversedStockIDs := []int64{}
		rows, err := tx.QueryContext(ctx, `SELECT id FROM reserved_stocks WHERE user_id = ? AND order_id = ?`, releaseStock.UserID, releaseStock.OrderID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			err := rows.Scan(&id)
			if err != nil {
				return err
			}
			reversedStockIDs = append(reversedStockIDs, id)
		}

		for _, reservedStockID := range reversedStockIDs {
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

func (r *WarehouseRepo) SetWarehouseStatus(ctx context.Context, warehouseID int64, isActive bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE warehouses SET is_active = ? WHERE id = ?`, isActive, warehouseID)
	if err != nil {
		return err
	}

	return nil
}

func (r *WarehouseRepo) TransferStockBetweenWarehouse(ctx context.Context, fromWarehouseID, toWarehouseID int64, productID string, quantity int64) error {
	err := dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		var err error
		var availableStock int64
		var shopID int64
		query := `SELECT quantity, shop_id FROM stocks WHERE product_id = ? AND warehouse_id = ?`
		err = tx.QueryRowContext(ctx, query, productID, fromWarehouseID).Scan(&availableStock, &shopID)
		if err != nil {
			return err
		}

		if availableStock < quantity {
			return fmt.Errorf("available stocks is less than request quantity")
		}

		qCheckWarehouse := "select is_active from warehouses where id = ?"
		var sourceWareHouseIsActive bool
		err = tx.QueryRowContext(ctx, qCheckWarehouse, fromWarehouseID).Scan(&sourceWareHouseIsActive)
		if err != nil {
			return err
		}

		if !sourceWareHouseIsActive {
			return errors.New("source warehouse is inactive")
		}

		var destinationWareHouseIsActive bool
		err = tx.QueryRowContext(ctx, qCheckWarehouse, toWarehouseID).Scan(&destinationWareHouseIsActive)
		if err != nil {
			return err
		}

		if !destinationWareHouseIsActive {
			return errors.New("destination warehouse is inactive")
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE 
				stocks
			SET
				quantity = quantity - ? 
			WHERE 
				product_id = ? 
			AND 
				warehouse_id = ?`,
			quantity, productID, fromWarehouseID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx,
			`INSERT INTO stocks (product_id, warehouse_id, shop_id, quantity)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(product_id, shop_id, warehouse_id)
		 DO UPDATE SET quantity = quantity + excluded.quantity`,
			productID, toWarehouseID, shopID, quantity)
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

func (pm *WarehouseRepo) GetWarehouseByIDs(ctx context.Context, productID ...uuid.UUID) ([]entity.Warehouse, error) {
	q := `select
		id,
		name,
		is_active
	from warehouses
	where id = ?`
	args := []any{}
	qMarks := buildPlaceHoldersInClause(len(productID))

	for _, v := range productID {
		args = append(args, v)
	}

	if len(productID) > 1 {
		q = `select
		id,
		name,
		is_active
	from warehouses
	where id IN (` + qMarks + `)`
	}
	rows, err := pm.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	warehouses := []entity.Warehouse{}

	for rows.Next() {
		var w entity.Warehouse
		err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.IsActive,
		)
		if err != nil {
			return nil, err
		}

		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func (pm *WarehouseRepo) GetWarehouseByShopID(ctx context.Context, shopID int64) ([]entity.Warehouse, error) {

	q := `
	SELECT w.id, w.name, w.is_active FROM stocks s 
	LEFT JOIN warehouses w ON w.id=s.warehouse_id 
	WHERE s.shop_id = ? GROUP BY w.id`

	rows, err := pm.db.QueryContext(ctx, q, shopID)
	if err != nil {
		return nil, err
	}

	warehouses := []entity.Warehouse{}

	for rows.Next() {
		var w entity.Warehouse
		err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.IsActive,
		)
		if err != nil {
			return nil, err
		}

		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func buildPlaceHoldersInClause(lenitems int) string {
	if lenitems == 0 {
		return ""
	}

	qMarks := strings.Repeat("?,", lenitems-1) + "?"
	return qMarks
}
