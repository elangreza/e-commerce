package service

import (
	"context"
	"errors"
	"github/elangreza/e-commerce/stock/internal/constanta"
	"github/elangreza/e-commerce/stock/internal/entity"
	"github/elangreza/e-commerce/stock/internal/params"

	"github.com/google/uuid"
)

type (
	stockRepo interface {
		GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error)
		ReserveStock(ctx context.Context, reserveStock entity.ReserveStock) ([]int64, error)
		ReleaseStock(ctx context.Context, releaseStock entity.ReleaseStock) ([]int64, error)
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

func (s *StockService) ReserveStock(ctx context.Context, reserveStock params.ReserveStock) ([]int64, error) {

	userID, ok := ctx.Value(constanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	stocks := make([]entity.Stock, len(reserveStock.Stocks))
	for i, stock := range reserveStock.Stocks {
		stocks[i] = entity.Stock{
			ProductID: stock.ProductID,
			Quantity:  stock.Quantity,
		}
	}

	reservedStockIDs, err := s.repo.ReserveStock(ctx, entity.ReserveStock{
		Stocks: stocks,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return reservedStockIDs, nil
}

func (s *StockService) ReleaseStock(ctx context.Context, releaseStock params.ReleaseStock) ([]int64, error) {

	userID, ok := ctx.Value(constanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	releasedStockIDs, err := s.repo.ReleaseStock(ctx, entity.ReleaseStock{
		ReservedStockIDs: releaseStock.ReservedStockIDs,
		UserID:           userID,
	})
	if err != nil {
		return nil, err
	}

	return releasedStockIDs, nil
}
