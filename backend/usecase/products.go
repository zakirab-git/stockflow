package usecase

import (
	"context"

	"github.com/zakirab-git/stockflow/dto"
	"github.com/zakirab-git/stockflow/model"
	"github.com/zakirab-git/stockflow/repo"
)

type ProductRepo interface {
	GetProducts(c context.Context, filter repo.ProductFilter) (products []model.Product, total int, err error)
}

type ProductUsecase struct {
	repo ProductRepo
}

func NewProductUsecase(rp ProductRepo) *ProductUsecase {
	return &ProductUsecase{repo: rp}
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
	if total > 0 {
		res.TotalPages = (total + req.Limit - 1) / req.Limit
	}

	for _, p := range products {
		res.Products = append(res.Products, dto.ProductFromModel(p))
	}

	return res, nil
}
