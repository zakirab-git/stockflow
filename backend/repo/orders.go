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
		RETURNING id, customer, total, status, created_at
	`
	addOrderItemQuery = `
		INSERT INTO order_items (order_id, product_id, sku, qty, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`
	decrementStockQuery = `
		UPDATE products SET stock = stock - $2 WHERE id = $1
	`
	getOrderByIDQuery = `
		SELECT id, customer, total, status, created_at
		FROM orders
		WHERE id = $1
	`
	getOrderByIDForUpdateQuery = `
		SELECT id, customer, total, status, created_at
		FROM orders
		WHERE id = $1
		FOR UPDATE
	`
	getOrderItemsQuery = `
		SELECT id, order_id, product_id, sku, qty, unit_price
		FROM order_items
		WHERE order_id = $1
	`
	updateOrderStatusQuery = `
		UPDATE orders SET status = $2 WHERE id = $1
	`
	restoreStockQuery = `
		UPDATE products SET stock = stock + $2 WHERE id = $1
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

func (r *OrderRepo) GetOrderByID(c context.Context, id int) (o model.Order, err error) {
	err = r.db.GetContext(c, &o, getOrderByIDQuery, id)
	return o, err
}

func (r *OrderRepo) GetOrderByIDForUpdate(c context.Context, id int) (o model.Order, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.GetContext(c, &o, getOrderByIDForUpdateQuery, id)
	return o, err
}

func (r *OrderRepo) GetOrderItems(c context.Context, orderID int) (items []model.OrderItem, err error) {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	err = getter.SelectContext(c, &items, getOrderItemsQuery, orderID)
	return items, err
}

func (r *OrderRepo) UpdateOrderStatus(c context.Context, id int, status string) error {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	_, err := getter.ExecContext(c, updateOrderStatusQuery, id, status)
	return err
}
func (r *OrderRepo) RestoreStock(c context.Context, productID int, qty int) error {
	getter := utils.Getter(r.db)
	if tx := utils.ExtractTx(c); tx != nil {
		getter = tx
	}

	_, err := getter.ExecContext(c, restoreStockQuery, productID, qty)
	return err
}
