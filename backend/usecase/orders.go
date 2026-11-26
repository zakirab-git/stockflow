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
				ProductID: p.ID,
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
		res.CreatedAt = order.CreatedAt
		return nil
	})

	return res, conflictSKU, available, err
}
