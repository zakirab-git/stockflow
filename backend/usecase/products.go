package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/zakirab-git/stockflow/apperrors"
	"github.com/zakirab-git/stockflow/dto"
	"github.com/zakirab-git/stockflow/model"
	"github.com/zakirab-git/stockflow/repo"
	"github.com/zakirab-git/stockflow/utils"
)

type ProductRepo interface {
	GetProducts(c context.Context, filter repo.ProductFilter) (products []model.Product, total int, err error)
	SKUExists(c context.Context, sku string) (exists bool, err error)
	AddProduct(c context.Context, p model.Product) (created model.Product, err error)
	UpdateProduct(c context.Context, p model.Product) (updated model.Product, err error)
	DeleteProduct(c context.Context, id int) error
	GetStockForUpdate(c context.Context, productID int) (stock int, err error)
	AddMovement(c context.Context, m model.StockMovement) (created model.StockMovement, err error)
	UpdateStock(c context.Context, productID int, stock int) error
	GetMovements(c context.Context, productID int) (movements []model.StockMovement, err error)
}

type ProductUsecase struct {
	repo ProductRepo
	txM  utils.TransactionxManagerInterface
}

func NewProductUsecase(rp ProductRepo, txM utils.TransactionxManagerInterface) *ProductUsecase {
	return &ProductUsecase{repo: rp, txM: txM}
}

func (uc *ProductUsecase) GetProducts(c context.Context, req dto.ProductListRequest) (res dto.ProductListResponse, err error) {
	products, total, err := uc.repo.GetProducts(c, repo.ProductFilter{
		Search: req.Search,
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	})
	if err != nil {
		return dto.ProductListResponse{}, err
	}

	res.Page = req.Page
	res.Limit = req.Limit
	res.Total = total
	res.Data = make([]dto.ProductResponse, 0, len(products))
	if total > 0 {
		res.TotalPages = (total + req.Limit - 1) / req.Limit
	}

	for _, p := range products {
		res.Data = append(res.Data, dto.ProductFromModel(p))
	}

	return res, nil
}

func (uc *ProductUsecase) AddProduct(c context.Context, req dto.AddProductRequest) (res dto.ProductResponse, err error) {
	exists, err := uc.repo.SKUExists(c, req.SKU)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	if exists {
		return dto.ProductResponse{}, errors.New("sku already exists")
	}

	created, err := uc.repo.AddProduct(c, model.Product{
		SKU:      req.SKU,
		Name:     req.Name,
		Price:    req.Price,
		Stock:    req.Stock,
		IsActive: req.IsActive,
	})
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ProductFromModel(created), nil
}

func (uc *ProductUsecase) UpdateProduct(c context.Context, id int, req dto.UpdateProductRequest) (res dto.ProductResponse, err error) {
	updated, err := uc.repo.UpdateProduct(c, model.Product{
		ID:       id,
		Name:     req.Name,
		Price:    req.Price,
		IsActive: req.IsActive,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return dto.ProductResponse{}, errors.New("product not found")
	}
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return dto.ProductFromModel(updated), nil
}

func (uc *ProductUsecase) DeleteProduct(c context.Context, id int) error {
	err := uc.repo.DeleteProduct(c, id)
	if err != nil && strings.Contains(err.Error(), "SQLSTATE 23503") {
		return errors.New("PRODUCT_HAS_HISTORY")
	}

	return err
}

func (uc *ProductUsecase) UpdateStock(c context.Context, productID int, req dto.UpdateStockRequest) (res dto.MovementResponse, available int, err error) {
	err = uc.txM.WithTxM(c, func(ctx context.Context) error {
		// Get stock
		stock, err := uc.repo.GetStockForUpdate(ctx, productID)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product not found")
		}
		if err != nil {
			return err
		}

		//Check if overdraw
		newStock := stock + req.Qty
		if req.Type == "out" {
			newStock = stock - req.Qty
			if newStock < 0 {
				available = stock
				return apperrors.ErrInsufficientStock
			}
		}

		movement, err := uc.repo.AddMovement(ctx, model.StockMovement{
			ProductID: productID,
			Type:      req.Type,
			Qty:       req.Qty,
			Note:      req.Note,
		})
		if err != nil {
			return err
		}

		if err := uc.repo.UpdateStock(ctx, productID, newStock); err != nil {
			return err
		}

		res = dto.MovementFromModel(movement)
		return nil
	})

	return res, available, err
}

func (uc *ProductUsecase) GetProductMovement(c context.Context, productID int) (res []dto.MovementResponse, err error) {
	movements, err := uc.repo.GetMovements(c, productID)
	if err != nil {
		return nil, err
	}

	for _, m := range movements {
		res = append(res, dto.MovementFromModel(m))
	}

	return res, nil
}
