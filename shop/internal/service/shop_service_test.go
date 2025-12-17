package service_test

import (
	"context"
	"testing"

	"github.com/elangreza/e-commerce/gen"
	globalcontanta "github.com/elangreza/e-commerce/pkg/globalcontanta"
	"github.com/elangreza/e-commerce/shop/internal/entity"
	"github.com/elangreza/e-commerce/shop/internal/service"
	"github.com/elangreza/e-commerce/shop/internal/service/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

type OrderServiceTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	mockWarehouseClient *mock.MockWarehouseServiceClient
	svc                 *service.ShopService
	mockShopRepo        *mock.MockShopRepo
}

func (s *OrderServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockWarehouseClient = mock.NewMockWarehouseServiceClient(s.ctrl)
	s.mockShopRepo = mock.NewMockShopRepo(s.ctrl)

	s.svc = service.NewShopService(
		s.mockShopRepo,
		s.mockWarehouseClient,
	)
}

func (s *OrderServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}

func (s *OrderServiceTestSuite) TestGetShops() {
	userID := uuid.New()

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.GetShopsRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.ShopList
	}{

		{
			name: "Success",
			req: &gen.GetShopsRequest{
				Ids:            []int64{1},
				WithWarehouses: true,
			},
			setupMock: func() {
				s.mockShopRepo.EXPECT().
					GetShopByIDs(gomock.Any(), gomock.Any()).
					Return([]entity.Shop{
						{
							ID:   1,
							Name: "test",
						},
					}, nil)

				s.mockWarehouseClient.EXPECT().
					GetWarehouseByShopID(gomock.Any(), gomock.Any()).
					Return(&gen.GetWarehouseByShopIDResponse{
						Warehouses: []*gen.Warehouse{
							{
								Id:       1,
								Name:     "test",
								IsActive: true,
							},
						},
					}, nil)

			},
			expectedError: "",
			expectedRes: &gen.ShopList{
				Shops: []*gen.Shop{
					{
						Id:   1,
						Name: "test",
						Warehouses: []*gen.Warehouse{
							{
								Id:       1,
								Name:     "test",
								IsActive: true,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetShops(ctx, tt.req)

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
