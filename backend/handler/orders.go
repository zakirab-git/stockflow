package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zakirab-git/stockflow/apperrors"
	"github.com/zakirab-git/stockflow/dto"
	"github.com/zakirab-git/stockflow/repo"
	"github.com/zakirab-git/stockflow/usecase"
)

type OrderUsecase interface {
	CreateOrder(c context.Context, req dto.CreateOrderRequest) (res dto.OrderResponse, conflictSKU string, available int, err error)
	GetOrder(c context.Context, id int) (res dto.OrderResponse, err error)
	CancelOrder(c context.Context, id int) (res dto.OrderResponse, err error)
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
		msg := "invalid request body"
		errText := err.Error()
		if strings.Contains(errText, "Items") && strings.Contains(errText, "min") {
			msg = "order must include at least one line item"
		} else if strings.Contains(errText, "Customer") {
			msg = "customer is required"
		}
		writeBadRequestMessage(c, msg)
		return
	}

	order, sku, available, err := h.usecase.CreateOrder(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, apperrors.ErrInsufficientStock) {
			writeAPIError(c, http.StatusConflict, dto.APIErrorBody{
				Code:      "INSUFFICIENT_STOCK",
				SKU:       sku,
				Available: &available,
			})
			return
		}
		if errors.Is(err, apperrors.ErrDuplicateSKU) {
			writeBadRequestMessage(c, "duplicate sku in order items")
			return
		}
		if unknownSKU, ok := usecase.IsUnknownSKU(err); ok {
			writeBadRequestMessage(c, "unknown sku: "+unknownSKU)
			return
		}
		if err.Error() == "order must include at least one line item" ||
			strings.HasPrefix(err.Error(), "quantity for sku") {
			writeBadRequestMessage(c, err.Error())
			return
		}
		writeAPIError(c, http.StatusInternalServerError, dto.APIErrorBody{Code: "INTERNAL_ERROR", Message: "failed to create order"})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeBadRequestMessage(c, "invalid order id")
		return
	}

	order, err := h.usecase.GetOrder(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "order not found" {
			writeAPIError(c, http.StatusNotFound, dto.APIErrorBody{Code: "NOT_FOUND", Message: "order not found"})
			return
		}
		writeAPIError(c, http.StatusInternalServerError, dto.APIErrorBody{Code: "INTERNAL_ERROR", Message: "failed to get order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeBadRequestMessage(c, "invalid order id")
		return
	}

	order, err := h.usecase.CancelOrder(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "ORDER_ALREADY_CANCELLED" {
			writeAPIError(c, http.StatusConflict, dto.APIErrorBody{Code: "ORDER_ALREADY_CANCELLED"})
			return
		}
		if err.Error() == "order not found" {
			writeAPIError(c, http.StatusNotFound, dto.APIErrorBody{Code: "NOT_FOUND", Message: "order not found"})
			return
		}
		if errors.Is(err, repo.ErrProductSKUNotFound) {
			writeAPIError(c, http.StatusConflict, dto.APIErrorBody{
				Code:    "PRODUCT_MISSING",
				Message: "cannot restore stock; product for an order line no longer exists",
			})
			return
		}
		writeAPIError(c, http.StatusInternalServerError, dto.APIErrorBody{Code: "INTERNAL_ERROR", Message: "failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, order)
}
