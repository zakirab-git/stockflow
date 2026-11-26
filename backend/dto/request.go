package dto

type ProductListRequest struct {
	Search *string `form:"search"`
	Page   int     `form:"page,default=1" binding:"min=1"`
	Limit  int     `form:"limit,default=10" binding:"min=1,max=20"`
}
