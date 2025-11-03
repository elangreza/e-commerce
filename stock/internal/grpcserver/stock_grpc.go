package grpcserver

import (
	"context"
	"github/elangreza/e-commerce/stock/internal/entity"
	"github/elangreza/e-commerce/stock/internal/params"

	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	stockService interface {
		GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error)
		ReserveStock(ctx context.Context, reserveStock params.ReserveStock) ([]int64, error)
		ReleaseStock(ctx context.Context, releaseStock params.ReleaseStock) ([]int64, error)
		ConfirmStock(ctx context.Context, confirmStock params.ConfirmStock) ([]int64, error)
	}

	StockGRPCServer struct {
		service stockService
		gen.UnimplementedStockServiceServer
	}
)

func NewStockGRPCServer(service stockService) *StockGRPCServer {
	return &StockGRPCServer{
		service: service,
	}
}

func (s *StockGRPCServer) GetStocks(ctx context.Context, req *gen.GetStockRequest) (*gen.StockList, error) {
	stocks, err := s.service.GetStocks(ctx, req.GetProductIds())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	stockList := &gen.StockList{}
	for _, stock := range stocks {
		stockList.Stocks = append(stockList.Stocks, &gen.Stock{
			ProductId: stock.ProductID.String(),
			Quantity:  stock.Quantity,
		})
	}

	return stockList, nil
}

func (s *StockGRPCServer) ReserveStock(ctx context.Context, req *gen.ReserveStockRequest) (*gen.ReserveStockResponse, error) {
	stocks := params.ReserveStock{}
	for _, stock := range req.GetStocks() {
		productID, err := uuid.Parse(stock.GetProductId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		stocks.Stocks = append(stocks.Stocks, params.Stock{
			ProductID: productID,
			Quantity:  stock.GetQuantity(),
		})
	}
	reservedStockIDs, err := s.service.ReserveStock(ctx, stocks)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.ReserveStockResponse{
		ReservedStockIds: reservedStockIDs,
	}, nil
}

func (s *StockGRPCServer) ReleaseStock(ctx context.Context, req *gen.ReleaseStockRequest) (*gen.ReleaseStockResponse, error) {
	releasedStockIDs, err := s.service.ReleaseStock(ctx, params.ReleaseStock{
		ReservedStockIDs: req.GetReservedStockIds(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.ReleaseStockResponse{
		ReleasedStockIds: releasedStockIDs,
	}, nil
}

func (s *StockGRPCServer) ConfirmedStock(ctx context.Context, req *gen.ConfirmedStockRequest) (*gen.ConfirmedStockResponse, error) {
	confirmedStockIDs, err := s.service.ConfirmStock(ctx, params.ConfirmStock{
		ReservedStockIDs: req.GetReservedStockIds(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.ConfirmedStockResponse{
		ConfirmedStockIds: confirmedStockIDs,
	}, nil
}
