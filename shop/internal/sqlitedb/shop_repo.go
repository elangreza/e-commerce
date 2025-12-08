package sqlitedb

import (
	"context"
	"database/sql"
	"strings"

	"github.com/elangreza/e-commerce/shop/internal/entity"
)

type ShopRepo struct {
	db *sql.DB
}

func NewShopRepo(db *sql.DB) *ShopRepo {
	return &ShopRepo{db: db}
}

func (pm *ShopRepo) GetShopByIDs(ctx context.Context, IDs ...int64) ([]entity.Shop, error) {
	q := `select
		id,
		name
	from shops
	where id = ?`
	args := []any{}
	qMarks := buildPlaceHoldersInClause(len(IDs))

	for _, v := range IDs {
		args = append(args, v)
	}

	if len(IDs) > 1 {
		q = `select
		id,
		name
	from shops
	where id IN (` + qMarks + `)`
	}
	rows, err := pm.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	shops := []entity.Shop{}

	for rows.Next() {
		var w entity.Shop
		err := rows.Scan(
			&w.ID,
			&w.Name,
		)
		if err != nil {
			return nil, err
		}

		shops = append(shops, w)
	}

	return shops, nil
}

func buildPlaceHoldersInClause(lenitems int) string {
	if lenitems == 0 {
		return ""
	}

	qMarks := strings.Repeat("?,", lenitems-1) + "?"
	return qMarks
}
