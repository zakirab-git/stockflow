package dto

import (
	"time"

	"github.com/zakirab-git/stockflow/model"
)

type ProductResponse struct {
	ID    int    `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
}

type MovementResponse struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Type      string    `json:"type"`
	Qty       int       `json:"qty"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemResponse struct {
	SKU       string `json:"sku"`
	Qty       int    `json:"qty"`
	UnitPrice int    `json:"unit_price"`
}

type OrderResponse struct {
	ID        int                 `json:"id"`
	Customer  string              `json:"customer"`
	Items     []OrderItemResponse `json:"items"`
	Total     int                 `json:"total"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
}

// -----------------------------------------------------------------------------
// Mapping functions (entity -> response DTO)
// -----------------------------------------------------------------------------

func ProductFromModel(p model.Product) ProductResponse {
	return ProductResponse{
		ID:    p.ID,
		SKU:   p.SKU,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}
}

func MovementFromModel(m model.StockMovement) MovementResponse {
	return MovementResponse{
		ID:        m.ID,
		ProductID: m.ProductID,
		Type:      m.Type,
		Qty:       m.Qty,
		Note:      m.Note,
		CreatedAt: m.CreatedAt,
	}
}
