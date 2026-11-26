package dto

type ProductListRequest struct {
	Search *string `form:"search"`
	Page   int     `form:"page,default=1" binding:"min=1"`
	Limit  int     `form:"limit,default=10" binding:"min=1,max=20"`
}

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
	Stock    int    `json:"stock" binding:"min=0"`
	IsActive bool   `json:"is_active"`
}
