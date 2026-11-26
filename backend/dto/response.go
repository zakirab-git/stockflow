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

func ProductFromModel(p model.Product) ProductResponse {
	return ProductResponse{
		ID:    p.ID,
		SKU:   p.SKU,
		Name:  p.Name,
		Price: p.Price,
		Stock: p.Stock,
	}
}

type MovementResponse struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Type      string    `json:"type"`
	Qty       int       `json:"qty"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
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
