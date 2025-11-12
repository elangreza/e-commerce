package sqlitedb

import (
	"context"
	"database/sql"
	"github/elangreza/e-commerce/pkg/dbsql"
	"github/elangreza/e-commerce/pkg/money"

	"github.com/elangreza/e-commerce/order/internal/entity"
	"github.com/google/uuid"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	// Implementation to retrieve cart by user ID from the database

	q := `
	SELECT id, user_id
	FROM carts WHERE user_id = ?;`

	var cart entity.Cart
	err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&cart.ID,
		&cart.UserID,
	)
	if err != nil {
		return nil, err
	}

	qItems := `
	SELECT 
		id, 
		cart_id, 
		product_id, 
		quantity, 
		price,
		currency
	FROM cart_items WHERE cart_id = ?;`

	rows, err := r.db.QueryContext(ctx, qItems, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cartItem entity.CartItem
		var price int64
		var currency string
		err = rows.Scan(
			&cartItem.ID,
			&cartItem.CartID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&price,
			&currency,
		)
		if err != nil {
			return nil, err
		}

		cartItem.Price, err = money.New(price, currency)
		if err != nil {
			return nil, err
		}

		cart.Items = append(cart.Items, cartItem)
	}

	return &cart, nil
}

func (r *CartRepository) CreateCart(ctx context.Context, cart entity.Cart) error {
	// Implementation to create a new cart in the database
	cartID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = dbsql.WithTransaction(r.db, func(tx *sql.Tx) error {
		q := `INSERT INTO carts (id, user_id) VALUES (?, ?);`
		_, err := tx.ExecContext(ctx, q,
			cartID,
			cart.UserID,
		)
		if err != nil {
			return err
		}

		qItem := `INSERT INTO cart_items 
		(id, cart_id, name, product_id, quantity, price, currency)
		VALUES (?,?,?,?,?,?,?);`

		for _, item := range cart.Items {

			cartItemID, err := uuid.NewV7()
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, qItem,
				cartItemID,
				cartID,
				item.Name,
				item.ProductID,
				item.Quantity,
				item.Price.Units,
				item.Price.CurrencyCode,
			)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepository) UpdateCartItem(ctx context.Context, item entity.CartItem) error {
	// Implementation to update an existing cart item in the database

	q := `UPDATE cart_items
		SET quantity = ?, name = ?, price = ?
		WHERE cart_id = ? AND product_id = ?;`
	_, err := r.db.ExecContext(ctx, q, item.Quantity, item.Name, item.Price, item.CartID, item.ProductID)
	if err != nil {
		return err
	}
	return nil
}
