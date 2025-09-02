package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
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
