package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"
	"errors"

	"github.com/elangreza/e-commerce/pkg/contextrequest"

	"github.com/elangreza/e-commerce/api/internal/constanta"
	params "github.com/elangreza/e-commerce/api/internal/params"
	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
)

func NewWarehouseService(pClient gen.WarehouseServiceClient) *WarehouseService {
	return &WarehouseService{
		WarehouseServiceClient: pClient,
	}
}

type WarehouseService struct {
	WarehouseServiceClient gen.WarehouseServiceClient
}

func (s *WarehouseService) SetWarehouseStatus(ctx context.Context, req params.SetWarehouseStatusRequest) error {
	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return errors.New("error when parsing userID")
	}

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

	_, err := s.WarehouseServiceClient.SetWarehouseStatus(newCtx, &gen.SetWarehouseStatusRequest{
		WarehouseId: req.WarehouseID,
		IsActive:    req.IsActive,
	})
	if err != nil {
		return convertErrGrpc(err)
	}

	return nil
}

func (s *WarehouseService) TransferStockBetweenWarehouse(ctx context.Context, req params.TransferStockBetweenWarehouseRequest) error {
	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return errors.New("error when parsing userID")
	}

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

	_, err := s.WarehouseServiceClient.TransferStockBetweenWarehouse(newCtx, &gen.TransferStockBetweenWarehouseRequest{
		FromWarehouseId: req.FromWarehouseId,
		ToWarehouseId:   req.ToWarehouseId,
		ProductId:       req.ProductId,
		Quantity:        req.Quantity,
	})
	if err != nil {
		return convertErrGrpc(err)
	}

	return nil
}
