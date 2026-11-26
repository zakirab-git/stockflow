package dto

// QUery params
type ProductListRequest struct {
	Search *string `form:"search"`
	Page   int     `form:"page,default=1" binding:"min=1"`
	Limit  int     `form:"limit,default=10" binding:"min=1,max=20"`
}

// Request body
type AddProductRequest struct {
	SKU      string `json:"sku" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Price    int    `json:"price" binding:"min=0"`
	Stock    int    `json:"stock" binding:"min=0"`
	IsActive bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    int    `json:"price" binding:"min=0"`
	IsActive bool   `json:"is_active"`
}

type UpdateStockRequest struct {
	Type string `json:"type" binding:"required,oneof=in out"`
	Qty  int    `json:"qty" binding:"required,min=1"`
	Note string `json:"note"`
}

type OrderItemRequest struct {
	SKU string `json:"sku" binding:"required"`
	Qty int    `json:"qty" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	Customer string             `json:"customer" binding:"required"`
	Items    []OrderItemRequest `json:"items" binding:"required,min=1"`
}
