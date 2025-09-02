package service

import (
	"context"
	"github/elangreza/e-commerce/stock/internal/entity"
)

type (
	stockRepo interface {
		GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error)
	}

	StockService struct {
		repo stockRepo
	}
)

func NewStockService(repo stockRepo) *StockService {
	return &StockService{
		repo: repo,
	}
}

func (s *StockService) GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error) {
	return s.repo.GetStocks(ctx, productIDs)
}
