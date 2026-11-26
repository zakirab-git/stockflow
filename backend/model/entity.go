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
