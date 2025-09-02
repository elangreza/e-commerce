package grpcserver

import (
	"context"
	"github/elangreza/e-commerce/stock/internal/entity"

	"github.com/elangreza/e-commerce/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	stockService interface {
		GetStocks(ctx context.Context, productIDs []string) ([]*entity.Stock, error)
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

func (s *StockGRPCServer) ListStocks(ctx context.Context, req *gen.ListStocksRequest) (*gen.ListStocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStocks not implemented")
}
func (s *StockGRPCServer) ReserveStock(ctx context.Context, req *gen.ReserveStockRequest) (*gen.ReserveStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReserveStock not implemented")
}
func (s *StockGRPCServer) ReleaseStock(ctx context.Context, req *gen.ReleaseStockRequest) (*gen.ReleaseStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseStock not implemented")
}
