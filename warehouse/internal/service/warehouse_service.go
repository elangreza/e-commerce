package service

import (
	"context"

	"github.com/elangreza/e-commerce/warehouse/internal/entity"

	"github.com/elangreza/e-commerce/pkg/extractor"

	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
)

type (
	warehouseRepo interface {
		GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error)
		ReserveStock(ctx context.Context, reserveStock entity.ReserveStock) ([]int64, error)
		ReleaseStock(ctx context.Context, releaseStock entity.ReleaseStock) ([]int64, error)
		SetWarehouseStatus(ctx context.Context, warehouseID int64, isActive bool) error
		TransferStockBetweenWarehouse(ctx context.Context, fromWarehouseID, toWarehouseID int64, productID string, quantity int64) error
		GetWarehouseByIDs(ctx context.Context, productID ...uuid.UUID) ([]entity.Warehouse, error)
		GetWarehouseByShopID(ctx context.Context, shopID int64) ([]entity.Warehouse, error)
	}

	WarehouseService struct {
		repo warehouseRepo
		gen.UnimplementedWarehouseServiceServer
	}
)

func NewWarehouseService(repo warehouseRepo) *WarehouseService {
	return &WarehouseService{
		repo: repo,
	}
}

func (s *WarehouseService) GetStocks(ctx context.Context, req *gen.GetStockRequest) (*gen.StockList, error) {
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

func (s *WarehouseService) ReserveStock(ctx context.Context, req *gen.ReserveStockRequest) (*gen.ReserveStockResponse, error) {
	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
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
		Stocks:  stocks,
		UserID:  userID,
		OrderID: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	return &gen.ReserveStockResponse{
		ReservedStockIds: reservedStockIDs,
	}, nil
}

func (s *WarehouseService) ReleaseStock(ctx context.Context, req *gen.ReleaseStockRequest) (*gen.ReleaseStockResponse, error) {
	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	releasedStockIDs, err := s.repo.ReleaseStock(ctx, entity.ReleaseStock{
		OrderID: req.OrderId,
		UserID:  userID,
	})
	if err != nil {
		return nil, err
	}

	return &gen.ReleaseStockResponse{
		ReleasedStockIds: releasedStockIDs,
	}, nil
}

func (s *WarehouseService) SetWarehouseStatus(ctx context.Context, req *gen.SetWarehouseStatusRequest) (*gen.Empty, error) {
	err := s.repo.SetWarehouseStatus(ctx, req.WarehouseId, req.GetIsActive())
	if err != nil {
		return nil, err
	}

	return &gen.Empty{}, nil
}

func (s *WarehouseService) TransferStockBetweenWarehouse(ctx context.Context, req *gen.TransferStockBetweenWarehouseRequest) (*gen.Empty, error) {
	err := s.repo.TransferStockBetweenWarehouse(ctx, req.FromWarehouseId, req.ToWarehouseId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &gen.Empty{}, nil
}

func (s *WarehouseService) GetWarehouseByShopID(ctx context.Context, req *gen.GetWarehouseByShopIDRequest) (*gen.GetWarehouseByShopIDResponse, error) {

	warehouses, err := s.repo.GetWarehouseByShopID(ctx, req.ShopId)
	if err != nil {
		return nil, err
	}

	res := []*gen.Warehouse{}
	for _, warehouse := range warehouses {
		res = append(res, &gen.Warehouse{
			Id:       warehouse.ID,
			Name:     warehouse.Name,
			IsActive: warehouse.IsActive,
		})
	}

	return &gen.GetWarehouseByShopIDResponse{
		Warehouses: res,
	}, nil
}
