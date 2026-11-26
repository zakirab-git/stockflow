package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zakirab-git/stockflow/model"
	"github.com/zakirab-git/stockflow/utils"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

type ProductFilter struct {
	Search *string
	Limit  int
	Offset int
}

const (
	getProductsQuery = `
		SELECT
			p.id, p.sku, p.name, p.price, p.stock, p.is_active, p.created_at,
			sm.created_at AS last_movement_at
		FROM products p
		LEFT JOIN LATERAL (
			SELECT created_at
			FROM stock_movements sm
			WHERE sm.product_id = p.id
			ORDER BY sm.created_at DESC
			LIMIT 1
		) sm ON TRUE
		WHERE ($1::text IS NULL OR p.name ILIKE '%' || $1 || '%' OR p.sku ILIKE '%' || $1 || '%')
		ORDER BY p.id
		LIMIT $2 OFFSET $3
	`
	countProductsQuery = `
		SELECT COUNT(*)
		FROM products p
		WHERE ($1::text IS NULL OR p.name ILIKE '%' || $1 || '%' OR p.sku ILIKE '%' || $1 || '%')
	`
)

func (r *ProductRepo) GetProducts(c context.Context, filter ProductFilter) (products []model.Product, total int, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	if err = getter.SelectContext(c, &products, getProductsQuery, filter.Search, filter.Limit, filter.Offset); err != nil {
		return nil, 0, err
	}

	if err = getter.GetContext(c, &total, countProductsQuery, filter.Search); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
