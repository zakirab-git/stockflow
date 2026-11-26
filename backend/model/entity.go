package model

import "time"

type Product struct {
	ID             int       `db:"id" json:"id"`
	SKU            string    `db:"sku" json:"sku"`
	Name           string    `db:"name" json:"name"`
	Price          int       `db:"price" json:"price"`
	Stock          int       `db:"stock" json:"stock"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	LastMovementAt *time.Time `db:"last_movement_at" json:"last_movement_at"`
}

type StockMovement struct {
	ID        int       `db:"id" json:"id"`
	ProductID int       `db:"product_id" json:"product_id"`
	Type      string    `db:"type" json:"type"`
	Qty       int       `db:"qty" json:"qty"`
	Note      string    `db:"note" json:"note"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
