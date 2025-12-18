package service_test

import (
	"context"
	"testing"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/pkg/globalcontanta"
	"github.com/elangreza/e-commerce/warehouse/internal/entity"
	"github.com/elangreza/e-commerce/warehouse/internal/service"
	"github.com/elangreza/e-commerce/warehouse/internal/service/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

type WarehouseServiceTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	svc               *service.WarehouseService
	mockWarehouseRepo *mock.MockwarehouseRepo
}

func (s *WarehouseServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.mockWarehouseRepo = mock.NewMockwarehouseRepo(s.ctrl)

	s.svc = service.NewWarehouseService(
		s.mockWarehouseRepo,
	)
}

func (s *WarehouseServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestWarehouseServiceSuite(t *testing.T) {
	suite.Run(t, new(WarehouseServiceTestSuite))
}

func (s *WarehouseServiceTestSuite) TestGetStocks() {
	productID := uuid.New()

	tests := []struct {
		name          string
		req           *gen.GetStockRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.StockList
	}{
		{
			name: "Success",
			req: &gen.GetStockRequest{
				ProductIds: []string{"1"},
			},
			setupMock: func() {
				s.mockWarehouseRepo.EXPECT().
					GetStocks(gomock.Any(), gomock.Any()).
					Return([]*entity.Stock{
						{
							ProductID: productID,
							Quantity:  10,
						},
					}, nil)
			},
			expectedError: "",
			expectedRes: &gen.StockList{
				Stocks: []*gen.Stock{
					{
						ProductId: productID.String(),
						Quantity:  10,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetStocks(context.Background(), tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(resp, tt.expectedRes)
			}
		})
	}
}

func (s *WarehouseServiceTestSuite) TestReserveStock() {
	userID := uuid.New()
	productID := uuid.New()
	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.ReserveStockRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.ReserveStockResponse
	}{
		{
			name: "Success",
			req: &gen.ReserveStockRequest{
				Stocks: []*gen.Stock{
					{
						ProductId: productID.String(),
						Quantity:  10,
					},
				},
				OrderId: "1",
			},
			setupMock: func() {
				s.mockWarehouseRepo.EXPECT().
					ReserveStock(gomock.Any(), gomock.Any()).
					Return([]int64{1}, nil)
			},
			expectedError: "",
			expectedRes: &gen.ReserveStockResponse{
				ReservedStockIds: []int64{1},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.ReserveStock(ctx, tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(resp, tt.expectedRes)
			}
		})
	}
}

func (s *WarehouseServiceTestSuite) TestReleaseStock() {
	userID := uuid.New()
	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.ReleaseStockRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.ReleaseStockResponse
	}{
		{
			name: "Success",
			req: &gen.ReleaseStockRequest{
				OrderId: "1",
			},
			setupMock: func() {
				s.mockWarehouseRepo.EXPECT().
					ReleaseStock(gomock.Any(), gomock.Any()).
					Return([]int64{1}, nil)
			},
			expectedError: "",
			expectedRes: &gen.ReleaseStockResponse{
				ReleasedStockIds: []int64{1},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.ReleaseStock(ctx, tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(resp, tt.expectedRes)
			}
		})
	}
}

func (s *WarehouseServiceTestSuite) TestSetWarehouseStatus() {
	userID := uuid.New()
	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.SetWarehouseStatusRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Success",
			req: &gen.SetWarehouseStatusRequest{
				WarehouseId: 1,
				IsActive:    false,
			},
			setupMock: func() {
				s.mockWarehouseRepo.EXPECT().
					SetWarehouseStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.SetWarehouseStatus(ctx, tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
			}
		})
	}
}

func (s *WarehouseServiceTestSuite) TestTransferStockBetweenWarehouse() {
	userID := uuid.New()
	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.TransferStockBetweenWarehouseRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Success",
			req: &gen.TransferStockBetweenWarehouseRequest{
				FromWarehouseId: 1,
				ToWarehouseId:   2,
				ProductId:       "1",
				Quantity:        10,
			},
			setupMock: func() {
				s.mockWarehouseRepo.EXPECT().
					TransferStockBetweenWarehouse(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.TransferStockBetweenWarehouse(ctx, tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
			}
		})
	}
}

func (s *WarehouseServiceTestSuite) TestGetWarehouseByShopID() {
	userID := uuid.New()
	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.GetWarehouseByShopIDRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.GetWarehouseByShopIDResponse
	}{
		{
			name: "Success",
			req: &gen.GetWarehouseByShopIDRequest{
				ShopId: 1,
			},
			setupMock: func() {
				s.mockWarehouseRepo.EXPECT().
					GetWarehouseByShopID(gomock.Any(), gomock.Any()).
					Return([]entity.Warehouse{
						{
							ID:       1,
							Name:     "a",
							IsActive: true,
						},
					}, nil)
			},
			expectedError: "",
			expectedRes: &gen.GetWarehouseByShopIDResponse{
				Warehouses: []*gen.Warehouse{
					{
						Id:       1,
						Name:     "a",
						IsActive: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetWarehouseByShopID(ctx, tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(resp, tt.expectedRes)
			}
		})
	}
}
