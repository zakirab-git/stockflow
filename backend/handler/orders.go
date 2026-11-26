package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/dto"
)

type OrderUsecase interface {
	CreateOrder(c context.Context, req dto.CreateOrderRequest) (res dto.OrderResponse, conflictSKU string, available int, err error)
}

type OrderHandler struct {
	usecase OrderUsecase
}

func NewOrderHandler(uc OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, sku, available, err := h.usecase.CreateOrder(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "INSUFFICIENT_STOCK" {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "INSUFFICIENT_STOCK", "sku": sku, "available": available}})
			return
		}
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}
