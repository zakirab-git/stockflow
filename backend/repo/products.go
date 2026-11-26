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
	skuExistsQuery = `
		SELECT EXISTS(SELECT 1 FROM products WHERE sku = $1)
	`
	addProductQuery = `
		INSERT INTO products (sku, name, price, stock, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sku, name, price, stock, is_active, created_at
	`
	updateProductQuery = `
		UPDATE products
		SET name = $2, price = $3, stock = $4, is_active = $5
		WHERE id = $1
		RETURNING id, sku, name, price, stock, is_active, created_at
	`
	deleteProductQuery = `
		DELETE FROM products WHERE id = $1
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

func (r *ProductRepo) SKUExists(c context.Context, sku string) (exists bool, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.GetContext(c, &exists, skuExistsQuery, sku)
	return exists, err
}

func (r *ProductRepo) AddProduct(c context.Context, p model.Product) (created model.Product, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.GetContext(c, &created, addProductQuery, p.SKU, p.Name, p.Price, p.Stock, p.IsActive)
	return created, err
}

func (r *ProductRepo) UpdateProduct(c context.Context, p model.Product) (updated model.Product, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.GetContext(c, &updated, updateProductQuery, p.ID, p.Name, p.Price, p.Stock, p.IsActive)
	return updated, err
}

func (r *ProductRepo) DeleteProduct(c context.Context, id int) error {
	_, err := r.db.ExecContext(c, deleteProductQuery, id)
	return err
}
