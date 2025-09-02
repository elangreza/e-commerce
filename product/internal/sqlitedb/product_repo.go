package sqlitedb

import (
	"context"
	"database/sql"

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

	q := `select
		id,
		name,
		description,
		price,
		image_url,
		created_at,
		updated_at
	from products
	where
		(name LIKE '%' || ? || '%' OR ? IS NULL)
	order by ` + req.OrderClause + ` LIMIT ? OFFSET ?`
	offset := (req.Page - 1) * req.Limit

	rows, err := pm.db.QueryContext(ctx, q, req.Search, req.Search, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageUrl, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (pm *ProductRepository) TotalProducts(ctx context.Context, req entity.ListProductRequest) (int64, error) {
	q := `select count(*) from products
	where
		(name LIKE '%' || ? || '%' OR ? IS NULL)`
	var total int64
	if err := pm.db.QueryRowContext(ctx, q, req.Search, req.Search).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (pm *ProductRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	q := `select
		id,
		name,
		description,
		price,
		image_url,
		created_at,
		updated_at
	from products
	where id = ?`
	var p entity.Product
	if err := pm.db.QueryRowContext(ctx, q, productID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageUrl, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}
