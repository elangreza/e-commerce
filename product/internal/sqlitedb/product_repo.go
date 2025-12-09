package sqlitedb

import (
	"context"
	"database/sql"
	"strings"

	"github.com/elangreza/e-commerce/pkg/money"

	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/google/uuid"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (pm *ProductRepository) ListProducts(ctx context.Context, req entity.ListProductRequest) ([]entity.Product, error) {
	// Start building WHERE conditions and args
	whereClauses := []string{"1=1"} // dummy condition to simplify logic
	args := []any{}

	// Name filter
	if req.Search != "" {
		whereClauses = append(whereClauses, "(name LIKE '%' || ? || '%')")
		args = append(args, req.Search)
	}

	// Build ORDER clause
	orderClause := ""
	if req.OrderClause != "" {
		orderClause = " ORDER BY " + req.OrderClause
	}

	// Build final query
	query := `SELECT id, name, description, price, currency, image_url, created_at, updated_at, shop_id
              FROM products
              WHERE ` + strings.Join(whereClauses, " AND ") + orderClause + ` LIMIT ? OFFSET ?`

	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)

	rows, err := pm.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		var priceAmount int64
		var priceCurrency string
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&priceAmount,
			&priceCurrency,
			&p.ImageUrl,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.ShopID,
		); err != nil {
			return nil, err
		}

		p.Price, err = money.New(priceAmount, priceCurrency)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (pm *ProductRepository) TotalProducts(ctx context.Context, req entity.ListProductRequest) (int64, error) {
	whereClauses := []string{"1=1"}
	args := []any{}

	if req.Search != "" {
		whereClauses = append(whereClauses, "(name LIKE '%' || ? || '%')")
		args = append(args, req.Search)
	}

	query := `SELECT COUNT(1) FROM products WHERE ` + strings.Join(whereClauses, " AND ")

	var total int64
	if err := pm.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (pm *ProductRepository) GetProductByIDs(ctx context.Context, productID ...uuid.UUID) ([]entity.Product, error) {
	q := `select
		id,
		name,
		description,
		price,
		currency,
		image_url,
		created_at,
		updated_at,
		shop_id
	from products
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
		description,
		price,
		currency,
		image_url,
		created_at,
		updated_at,
		shop_id
	from products
	where id IN (` + qMarks + `)`
	}
	rows, err := pm.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	products := []entity.Product{}

	for rows.Next() {
		var p entity.Product
		var priceAmount int64
		var priceCurrency string
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&priceAmount,
			&priceCurrency,
			&p.ImageUrl,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.ShopID)
		if err != nil {
			return nil, err
		}
		p.Price, err = money.New(priceAmount, priceCurrency)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func buildPlaceHoldersInClause(lenitems int) string {
	if lenitems == 0 {
		return ""
	}

	qMarks := strings.Repeat("?,", lenitems-1) + "?"
	return qMarks
}
