package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zakirab-git/stockflow/dto"
	"github.com/zakirab-git/stockflow/model"
	"github.com/zakirab-git/stockflow/utils"
)

type OrderRepo interface {
	GetProductBySKUForUpdate(c context.Context, sku string) (p model.Product, err error)
	CreateOrder(c context.Context, o model.Order) (created model.Order, err error)
	AddOrderItem(c context.Context, item model.OrderItem) error
	DecrementStock(c context.Context, productID int, qty int) error
	GetOrderByID(c context.Context, id int) (o model.Order, err error)
	GetOrderByIDForUpdate(c context.Context, id int) (o model.Order, err error)
	GetOrderItems(c context.Context, orderID int) (items []model.OrderItem, err error)
	CancelOrder(c context.Context, id int) error
	RestoreStock(c context.Context, sku string, qty int) error
}

type OrderUsecase struct {
	repo OrderRepo
	txM  utils.TransactionxManagerInterface
}

func NewOrderUsecase(rp OrderRepo, txM utils.TransactionxManagerInterface) *OrderUsecase {
	return &OrderUsecase{repo: rp, txM: txM}
}

func (uc *OrderUsecase) CreateOrder(c context.Context, req dto.CreateOrderRequest) (res dto.OrderResponse, conflictSKU string, available int, err error) {
	err = uc.txM.WithTxM(c, func(ctx context.Context) error {
		var products []model.Product
		total := 0
		//check items & calculate total
		for _, item := range req.Items {
			p, err := uc.repo.GetProductBySKUForUpdate(ctx, item.SKU)
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New("product not found")
			}
			if err != nil {
				return err
			}
			if p.Stock < item.Qty {
				conflictSKU = item.SKU
				available = p.Stock
				return errors.New("INSUFFICIENT_STOCK")
			}

			total += p.Price * item.Qty
			products = append(products, p)
		}
		// insert to order table
		order, err := uc.repo.CreateOrder(ctx, model.Order{Customer: req.Customer, Total: total})
		if err != nil {
			return err
		}
		// decrement stock on each item
		for i, item := range req.Items {
			p := products[i]

			if err := uc.repo.DecrementStock(ctx, p.ID, item.Qty); err != nil {
				return err
			}

			if err := uc.repo.AddOrderItem(ctx, model.OrderItem{
				OrderID:   order.ID,
				SKU:       p.SKU,
				Qty:       item.Qty,
				UnitPrice: p.Price,
			}); err != nil {
				return err
			}

			res.Items = append(res.Items, dto.OrderItemResponse{
				SKU:       p.SKU,
				Qty:       item.Qty,
				UnitPrice: p.Price,
			})
		}

		res.ID = order.ID
		res.Customer = order.Customer
		res.Total = order.Total
		res.Status = order.Status
		res.CreatedAt = order.CreatedAt
		return nil
	})

	return res, conflictSKU, available, err
}

func (uc *OrderUsecase) GetOrder(c context.Context, id int) (res dto.OrderResponse, err error) {
	order, err := uc.repo.GetOrderByID(c, id)
	if errors.Is(err, sql.ErrNoRows) {
		return res, errors.New("order not found")
	}
	if err != nil {
		return res, err
	}

	items, err := uc.repo.GetOrderItems(c, id)
	if err != nil {
		return res, err
	}

	res = orderToResponse(order)
	for _, item := range items {
		res.Items = append(res.Items, orderItemToResponse(item))
	}

	return res, nil
}

func (uc *OrderUsecase) CancelOrder(c context.Context, id int) (res dto.OrderResponse, err error) {
	err = uc.txM.WithTxM(c, func(ctx context.Context) error {
		order, err := uc.repo.GetOrderByIDForUpdate(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("order not found")
		}
		if err != nil {
			return err
		}
		if order.Status == "CANCELLED" {
			return errors.New("ORDER_ALREADY_CANCELLED")
		}

		items, err := uc.repo.GetOrderItems(ctx, id)
		if err != nil {
			return err
		}

		for _, item := range items {
			if err := uc.repo.RestoreStock(ctx, item.SKU, item.Qty); err != nil {
				return err
			}
		}

		if err := uc.repo.CancelOrder(ctx, id); err != nil {
			return err
		}

		order.Status = "CANCELLED"
		res = orderToResponse(order)
		for _, item := range items {
			res.Items = append(res.Items, orderItemToResponse(item))
		}

		return nil
	})

	return res, err
}

//mappings

func orderToResponse(o model.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:        o.ID,
		Customer:  o.Customer,
		Total:     o.Total,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
	}
}

func orderItemToResponse(item model.OrderItem) dto.OrderItemResponse {
	return dto.OrderItemResponse{
		SKU:       item.SKU,
		Qty:       item.Qty,
		UnitPrice: item.UnitPrice,
	}
}
