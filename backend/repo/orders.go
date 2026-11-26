package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zakirab-git/stockflow/model"
	"github.com/zakirab-git/stockflow/utils"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

const (
	getProductBySKUForUpdateQuery = `
		SELECT id, sku, name, price, stock, is_active, created_at
		FROM products
		WHERE sku = $1
		FOR UPDATE
	`
	createOrderQuery = `
		INSERT INTO orders (customer, total)
		VALUES ($1, $2)
		RETURNING id, customer, total, created_at
	`
	addOrderItemQuery = `
		INSERT INTO order_items (order_id, product_id, sku, qty, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`
	decrementStockQuery = `
		UPDATE products SET stock = stock - $2 WHERE id = $1
	`
)

func (r *OrderRepo) GetProductBySKUForUpdate(c context.Context, sku string) (p model.Product, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.GetContext(c, &p, getProductBySKUForUpdateQuery, sku)
	return p, err
}

func (r *OrderRepo) CreateOrder(c context.Context, o model.Order) (created model.Order, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.GetContext(c, &created, createOrderQuery, o.Customer, o.Total)
	return created, err
}

func (r *OrderRepo) AddOrderItem(c context.Context, item model.OrderItem) error {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	_, err := getter.ExecContext(c, addOrderItemQuery, item.OrderID, item.ProductID, item.SKU, item.Qty, item.UnitPrice)
	return err
}

func (r *OrderRepo) DecrementStock(c context.Context, productID int, qty int) error {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}
	_, err := getter.ExecContext(c, decrementStockQuery, productID, qty)
	return err
}
