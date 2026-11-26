package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/dto"
)

type ProductUsecase interface {
	GetProducts(c context.Context, data dto.ProductListRequest) (res dto.ProductListResponse, err error)
	AddProduct(c context.Context, req dto.AddProductRequest) (res dto.ProductResponse, err error)
	UpdateProduct(c context.Context, id int, req dto.UpdateProductRequest) (res dto.ProductResponse, err error)
	DeleteProduct(c context.Context, id int) error
	UpdateStock(c context.Context, id int, req dto.UpdateStockRequest) (res dto.MovementResponse, available int, err error)
	GetProductMovement(c context.Context, id int) (res []dto.MovementResponse, err error)
}

type Handlers struct {
	Product *ProductHandler
	Order   *OrderHandler
}

func NewHandlers(product *ProductHandler, order *OrderHandler) *Handlers {
	return &Handlers{Product: product, Order: order}
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

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var req dto.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.usecase.AddProduct(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "sku already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.usecase.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := h.usecase.DeleteProduct(c.Request.Context(), id); err != nil {
		if err.Error() == "PRODUCT_HAS_HISTORY" {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "PRODUCT_HAS_HISTORY"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	var req dto.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movement, available, err := h.usecase.UpdateStock(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "INSUFFICIENT_STOCK" {
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "INSUFFICIENT_STOCK", "available": available}})
			return
		}
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, movement)
}

func (h *ProductHandler) GetProductMovement(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	movements, err := h.usecase.GetProductMovement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movements})
}
