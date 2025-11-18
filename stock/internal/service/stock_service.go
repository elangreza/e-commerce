package service

import (
	"context"
	"errors"
	"github/elangreza/e-commerce/stock/internal/constanta"
	"github/elangreza/e-commerce/stock/internal/entity"

	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
)

type (
	stockRepo interface {
		GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error)
		ReserveStock(ctx context.Context, reserveStock entity.ReserveStock) ([]int64, error)
		ReleaseStock(ctx context.Context, releaseStock entity.ReleaseStock) ([]int64, error)
		ConfirmStock(ctx context.Context, confirmStock entity.ConfirmStock) ([]int64, error)
	}

	StockService struct {
		repo stockRepo
		gen.UnimplementedStockServiceServer
	}
)

func NewStockService(repo stockRepo) *StockService {
	return &StockService{
		repo: repo,
	}
}

func (s *StockService) GetStocks(ctx context.Context, req *gen.GetStockRequest) (*gen.StockList, error) {
	stocks, err := s.repo.GetStocks(ctx, req.ProductIds)
	if err != nil {
		return nil, err
	}
	res := []*gen.Stock{}
	for _, stock := range stocks {
		res = append(res, &gen.Stock{
			ProductId: stock.ProductID.String(),
			Quantity:  stock.Quantity,
		})
	}
	return &gen.StockList{
		Stocks: res,
	}, nil
}

func (s *StockService) ReserveStock(ctx context.Context, req *gen.ReserveStockRequest) (*gen.ReserveStockResponse, error) {
	userID, ok := ctx.Value(constanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	stocks := make([]entity.Stock, len(req.Stocks))
	for i, stock := range req.Stocks {
		productID, err := uuid.Parse(stock.ProductId)
		if err != nil {
			return nil, err
		}

		stocks[i] = entity.Stock{
			ProductID: productID,
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

	return &gen.ReserveStockResponse{
		ReservedStockIds: reservedStockIDs,
	}, nil
}

func (s *StockService) ReleaseStock(ctx context.Context, req *gen.ReleaseStockRequest) (*gen.ReleaseStockResponse, error) {
	userID, ok := ctx.Value(constanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	releasedStockIDs, err := s.repo.ReleaseStock(ctx, entity.ReleaseStock{
		ReservedStockIDs: req.ReservedStockIds,
		UserID:           userID,
	})
	if err != nil {
		return nil, err
	}

	return &gen.ReleaseStockResponse{
		ReleasedStockIds: releasedStockIDs,
	}, nil
}

func (s *StockService) ConfirmedStock(ctx context.Context, req *gen.ConfirmedStockRequest) (*gen.ConfirmedStockResponse, error) {

	userID, ok := ctx.Value(constanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	confirmedStockIDs, err := s.repo.ConfirmStock(ctx, entity.ConfirmStock{
		ReservedStockIDs: req.ReservedStockIds,
		UserID:           userID,
	})
	if err != nil {
		return nil, err
	}

	return &gen.ConfirmedStockResponse{
		ConfirmedStockIds: confirmedStockIDs,
	}, nil
}
