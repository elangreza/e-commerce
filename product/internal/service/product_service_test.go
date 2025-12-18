package service_test

import (
	"context"
	"testing"

	"github.com/elangreza/e-commerce/gen"
	globalcontanta "github.com/elangreza/e-commerce/pkg/globalcontanta"
	"github.com/elangreza/e-commerce/product/internal/entity"
	"github.com/elangreza/e-commerce/product/internal/service"
	"github.com/elangreza/e-commerce/product/internal/service/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

type ProductServiceTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	mockWarehouseClient *mock.MockWarehouseServiceClient
	svc                 *service.ProductService
	mockProductRepo     *mock.MockproductRepo
}

func (s *ProductServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockWarehouseClient = mock.NewMockWarehouseServiceClient(s.ctrl)
	s.mockProductRepo = mock.NewMockproductRepo(s.ctrl)

	s.svc = service.NewProductService(
		s.mockProductRepo,
		s.mockWarehouseClient,
	)
}

func (s *ProductServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) TestListProducts() {
	userID := uuid.New()

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	productID := uuid.New()

	tests := []struct {
		name          string
		req           *gen.ListProductsRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.ListProductsResponse
	}{

		{
			name: "Success",
			req: &gen.ListProductsRequest{
				Search:    "",
				Limit:     0,
				Page:      0,
				SortBy:    "",
				WithStock: true,
			},
			setupMock: func() {
				s.mockProductRepo.EXPECT().
					ListProducts(gomock.Any(), gomock.Any()).
					Return([]entity.Product{
						{
							ID:          productID,
							Name:        "",
							Description: "",
							Price:       &gen.Money{},
							ImageUrl:    "",
							CreatedAt:   "",
							UpdatedAt:   "",
							ShopID:      0,
						},
					}, nil)

				s.mockProductRepo.EXPECT().
					TotalProducts(gomock.Any(), gomock.Any()).
					Return(int64(1), nil)

				s.mockWarehouseClient.EXPECT().
					GetStocks(gomock.Any(), gomock.Any()).
					Return(&gen.StockList{
						Stocks: []*gen.Stock{
							{
								ProductId: productID.String(),
								Quantity:  1,
							},
						},
					}, nil)

			},
			expectedError: "",
			expectedRes: &gen.ListProductsResponse{
				Products: []*gen.Product{
					{
						Id:          productID.String(),
						Name:        "",
						Description: "",
						ImageUrl:    "",
						Price:       &gen.Money{},
						Stock:       1,
						ShopId:      0,
					},
				},
				Total:      1,
				TotalPages: 1,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.ListProducts(ctx, tt.req)

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

func (s *ProductServiceTestSuite) TestGetProducts() {
	userID := uuid.New()

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	productID := uuid.New()

	tests := []struct {
		name          string
		req           *gen.GetProductsRequest
		setupMock     func()
		expectedError string
		expectedRes   *gen.Products
	}{

		{
			name: "Success",
			req: &gen.GetProductsRequest{
				Ids:       []string{productID.String()},
				WithStock: true,
			},
			setupMock: func() {
				s.mockProductRepo.EXPECT().
					GetProductByIDs(gomock.Any(), gomock.Any()).
					Return([]entity.Product{
						{
							ID:          productID,
							Name:        "",
							Description: "",
							Price:       &gen.Money{},
							ImageUrl:    "",
							CreatedAt:   "",
							UpdatedAt:   "",
							ShopID:      0,
						},
					}, nil)

				s.mockWarehouseClient.EXPECT().
					GetStocks(gomock.Any(), gomock.Any()).
					Return(&gen.StockList{
						Stocks: []*gen.Stock{
							{
								ProductId: productID.String(),
								Quantity:  1,
							},
						},
					}, nil)

			},
			expectedError: "",
			expectedRes: &gen.Products{
				Products: []*gen.Product{
					{
						Id:          productID.String(),
						Name:        "",
						Description: "",
						ImageUrl:    "",
						Price:       &gen.Money{},
						Stock:       1,
						ShopId:      0,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetProducts(ctx, tt.req)

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
