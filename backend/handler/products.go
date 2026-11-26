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
