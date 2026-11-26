package dto

import "github.com/zakirab-git/stockflow/model"

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
