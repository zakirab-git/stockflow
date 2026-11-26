package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/dto"
)

type ProductUsecase interface {
	GetProducts(c context.Context, data dto.ProductListRequest) (res dto.ProductListResponse, err error)
}

type Handlers struct {
	Product *ProductHandler
}

func NewHandlers(product *ProductHandler) *Handlers {
	return &Handlers{Product: product}
}

type ProductHandler struct {
	usecase ProductUsecase
}

func NewProductHandler(uc ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: uc}
}

// APIs

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var req dto.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.usecase.GetProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
