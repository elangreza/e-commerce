package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/constanta"
	"github.com/elangreza/e-commerce/order/internal/entity"
	"github.com/elangreza/e-commerce/order/internal/service"
	"github.com/elangreza/e-commerce/order/internal/service/mock"
	globalcontanta "github.com/elangreza/e-commerce/pkg/globalcontanta"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

type OrderServiceTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	mockOrderRepo       *mock.MockorderRepo
	mockCartRepo        *mock.MockcartRepo
	mockWarehouseClient *mock.MockWarehouseServiceClient
	mockProductClient   *mock.MockProductServiceClient
	mockPaymentClient   *mock.MockPaymentServiceClient
	svc                 *service.OrderService
}

func (s *OrderServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockOrderRepo = mock.NewMockorderRepo(s.ctrl)
	s.mockCartRepo = mock.NewMockcartRepo(s.ctrl)
	s.mockWarehouseClient = mock.NewMockWarehouseServiceClient(s.ctrl)
	s.mockProductClient = mock.NewMockProductServiceClient(s.ctrl)
	s.mockPaymentClient = mock.NewMockPaymentServiceClient(s.ctrl)

	s.svc = service.NewOrderService(
		s.mockOrderRepo,
		s.mockCartRepo,
		s.mockWarehouseClient,
		s.mockProductClient,
		s.mockPaymentClient,
	)
}

func (s *OrderServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}

func (s *OrderServiceTestSuite) TestAddProductToCart() {
	userID := uuid.New()
	cartID := uuid.New()
	productID := "prod-123"

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.AddCartItemRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Failed because stock is not enough",
			req: &gen.AddCartItemRequest{
				ProductId: productID,
				Quantity:  2,
			},
			setupMock: func() {
				// 1. Get Cart
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(nil, sql.ErrNoRows)

				// 2. get product with stock
				s.mockProductClient.EXPECT().GetProducts(ctx, &gen.GetProductsRequest{
					Ids:       []string{productID},
					WithStock: true,
				}).Return(&gen.Products{
					Products: []*gen.Product{
						{
							Id:    productID,
							Name:  "Test Product",
							Stock: 1,
							Price: &gen.Money{Units: 10000, CurrencyCode: "IDR"},
						},
					},
				}, nil)
			},
			expectedError: "quantity cannot exceed the maximum stock, current stock is 1",
		},
		{
			name: "Success",
			req: &gen.AddCartItemRequest{
				ProductId: productID,
				Quantity:  2,
			},
			setupMock: func() {
				// 1. Get Cart
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(nil, sql.ErrNoRows)

				// 2. get product with stock
				s.mockProductClient.EXPECT().GetProducts(ctx, &gen.GetProductsRequest{
					Ids:       []string{productID},
					WithStock: true,
				}).Return(&gen.Products{
					Products: []*gen.Product{
						{
							Id:          productID,
							Name:        "a",
							Description: "a",
							ImageUrl:    "a",
							Price: &gen.Money{
								Units:        10000,
								CurrencyCode: "IDR",
							},
							Stock:  2,
							ShopId: 1,
						},
					},
				}, nil)

				// 3. Create Cart assuming all request validation is passed
				s.mockCartRepo.EXPECT().
					CreateCart(ctx,
						entity.Cart{
							UserID: userID,
							Items: []entity.CartItem{
								{
									ProductID: productID,
									Quantity:  2,
									Name:      "a",
									Price: &gen.Money{
										Units:        10000,
										CurrencyCode: "IDR",
									},
								},
							},
						},
					).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Success but cart already exists",
			req: &gen.AddCartItemRequest{
				ProductId: productID,
				Quantity:  2,
			},
			setupMock: func() {
				// 1. Get Cart
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(&entity.Cart{
						ID:     cartID,
						UserID: userID,
						Items: []entity.CartItem{
							{ProductID: productID, Quantity: 2},
						},
					}, nil)

				// 2. get product with stock
				s.mockProductClient.EXPECT().GetProducts(ctx, &gen.GetProductsRequest{
					Ids:       []string{productID},
					WithStock: true,
				}).Return(&gen.Products{
					Products: []*gen.Product{
						{
							Id:          productID,
							Name:        "a",
							Description: "a",
							ImageUrl:    "a",
							Price: &gen.Money{
								Units:        10000,
								CurrencyCode: "IDR",
							},
							Stock:  2,
							ShopId: 1,
						},
					},
				}, nil)

				// 3. Update Cart since the cart already exists
				s.mockCartRepo.EXPECT().
					UpdateCartItem(ctx,
						entity.CartItem{
							CartID:    cartID,
							ProductID: productID,
							Quantity:  2,
							Name:      "a",
							Price: &gen.Money{
								Units:        10000,
								CurrencyCode: "IDR",
							},
						},
					).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()
			_, err := s.svc.AddProductToCart(ctx, tt.req)
			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
			} else {
				s.NoError(err)
			}
		})
	}

}

func (s *OrderServiceTestSuite) TestGetCart() {
	userID := uuid.New()
	cartID := uuid.New()
	productID := "prod-123"

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		setupMock     func()
		expectedError string
		expectedResp  *gen.Cart
	}{
		{
			name: "Success",
			setupMock: func() {
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(&entity.Cart{
						ID:     cartID,
						UserID: userID,
						Items: []entity.CartItem{
							{
								ID:        cartID,
								CartID:    cartID,
								ProductID: productID,
								Quantity:  1,
								Name:      "a",
								Price: &gen.Money{
									Units:        10000,
									CurrencyCode: "IDR",
								},
							},
						},
					}, nil)

				s.mockWarehouseClient.EXPECT().
					GetStocks(gomock.Any(), &gen.GetStockRequest{
						ProductIds: []string{productID},
					}).
					Return(&gen.StockList{
						Stocks: []*gen.Stock{
							{
								ProductId: productID,
								Quantity:  2,
							},
						},
					}, nil)
			},
			expectedError: "",
			expectedResp: &gen.Cart{
				Id: cartID.String(),
				Items: []*gen.CartItem{
					{
						ProductId: productID,
						Quantity:  1,
						Name:      "a",
						Price: &gen.Money{
							Units:        10000,
							CurrencyCode: "IDR",
						},
						ActualStock: 2,
					},
				},
			},
		},
		{
			name: "Error cart not found",
			setupMock: func() {
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(nil, sql.ErrNoRows)
			},
			expectedError: "cart not found",
			expectedResp:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetCart(ctx, &gen.Empty{})
			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(tt.expectedResp, resp)
			}
		})
	}
}

func (s *OrderServiceTestSuite) TestCreateOrder() {
	userID := uuid.New()
	cartID := uuid.New()
	productID := "prod-123"
	idempotencyKey := uuid.New()
	transactionID := "txn-abc-123"

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.CreateOrderRequest
		setupMock     func()
		expectedError string
		expectedResp  *gen.Order
	}{
		{
			name: "Success",
			req: &gen.CreateOrderRequest{
				IdempotencyKey: idempotencyKey.String(),
			},
			setupMock: func() {
				// 1. Check Idempotency
				s.mockOrderRepo.EXPECT().
					GetOrderByIdempotencyKey(gomock.Any(), idempotencyKey).
					Return(nil, sql.ErrNoRows)

				// 2. Get Cart
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(&entity.Cart{
						ID:     cartID,
						UserID: userID,
						Items: []entity.CartItem{
							{ProductID: productID, Quantity: 2},
						},
					}, nil)

				// 3. Get Products
				s.mockProductClient.EXPECT().
					GetProducts(gomock.Any(), &gen.GetProductsRequest{
						Ids:       []string{productID},
						WithStock: false,
					}).
					Return(&gen.Products{
						Products: []*gen.Product{
							{
								Id:    productID,
								Name:  "Test Product",
								Stock: 100,
								Price: &gen.Money{Units: 10000, CurrencyCode: "IDR"},
							},
						},
					}, nil)

				// 4. Create Order
				orderID := uuid.New()
				s.mockOrderRepo.EXPECT().
					CreateOrder(gomock.Any(), gomock.AssignableToTypeOf(entity.Order{})).
					Return(orderID, nil)

				// 5. Reserve Stock
				s.mockWarehouseClient.EXPECT().
					ReserveStock(gomock.Any(), gomock.AssignableToTypeOf(&gen.ReserveStockRequest{})).
					Return(&gen.ReserveStockResponse{ReservedStockIds: []int64{1}}, nil)

				// 6. Process Payment
				s.mockPaymentClient.EXPECT().
					ProcessPayment(gomock.Any(), gomock.AssignableToTypeOf(&gen.ProcessPaymentRequest{})).
					Return(&gen.ProcessPaymentResponse{
						TransactionId: transactionID,
					}, nil)

				// 7. Update Order Status
				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(2), orderID).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusStockReserved, payloads["status"])
						s.Equal(transactionID, payloads["transaction_id"])
						return nil
					})
			},
			expectedResp: &gen.Order{
				Status:        constanta.OrderStatusStockReserved.String(),
				TransactionId: transactionID,
			},
		},
		{
			name: "Failed_StockReservation",
			req: &gen.CreateOrderRequest{
				IdempotencyKey: idempotencyKey.String(),
			},
			setupMock: func() {
				// 1. Check Idempotency
				s.mockOrderRepo.EXPECT().
					GetOrderByIdempotencyKey(gomock.Any(), idempotencyKey).
					Return(nil, sql.ErrNoRows)

				// 2. Get Cart
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(&entity.Cart{
						ID:     cartID,
						UserID: userID,
						Items: []entity.CartItem{
							{ProductID: productID, Quantity: 2},
						},
					}, nil)

				// 3. Get Products
				s.mockProductClient.EXPECT().
					GetProducts(gomock.Any(), gomock.Any()).
					Return(&gen.Products{
						Products: []*gen.Product{
							{
								Id:    productID,
								Name:  "Test Product",
								Stock: 100,
								Price: &gen.Money{Units: 10000, CurrencyCode: "IDR"},
							},
						},
					}, nil)

				// 4. Create Order
				orderID := uuid.New()
				s.mockOrderRepo.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(orderID, nil)

				// 5. Reserve Stock FAILS
				s.mockWarehouseClient.EXPECT().
					ReserveStock(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("insufficient stock"))

				// 6. Rollback (Update Order to Failed)
				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(1), orderID).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusFailed, payloads["status"])
						return nil
					})
			},
			expectedError: "failed to reserve stock",
		},
		{
			name: "Failed_PaymentProcessing",
			req: &gen.CreateOrderRequest{
				IdempotencyKey: idempotencyKey.String(),
			},
			setupMock: func() {
				// 1. Check Idempotency
				s.mockOrderRepo.EXPECT().
					GetOrderByIdempotencyKey(gomock.Any(), idempotencyKey).
					Return(nil, sql.ErrNoRows)

				// 2. Get Cart
				s.mockCartRepo.EXPECT().
					GetCartByUserID(gomock.Any(), userID).
					Return(&entity.Cart{
						ID:     cartID,
						UserID: userID,
						Items: []entity.CartItem{
							{ProductID: productID, Quantity: 2},
						},
					}, nil)

				// 3. Get Products
				s.mockProductClient.EXPECT().
					GetProducts(gomock.Any(), gomock.Any()).
					Return(&gen.Products{
						Products: []*gen.Product{
							{
								Id:    productID,
								Name:  "Test Product",
								Stock: 100,
								Price: &gen.Money{Units: 10000, CurrencyCode: "IDR"},
							},
						},
					}, nil)

				// 4. Create Order
				orderID := uuid.New()
				s.mockOrderRepo.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(orderID, nil)

				// 5. Reserve Stock SUCCESS
				s.mockWarehouseClient.EXPECT().
					ReserveStock(gomock.Any(), gomock.Any()).
					Return(&gen.ReserveStockResponse{
						ReservedStockIds: []int64{1},
					}, nil)

				// 6. Process Payment FAILS
				s.mockPaymentClient.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("failed to process payment"))

				// 7. Release Stock
				s.mockWarehouseClient.EXPECT().
					ReleaseStock(gomock.Any(), gomock.Any()).
					Return(&gen.ReleaseStockResponse{
						ReleasedStockIds: []int64{1},
					}, nil)

				// 8. Update Order to Failed
				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(1), orderID).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusFailed, payloads["status"])
						return nil
					})
			},
			expectedError: "failed to process payment",
			expectedResp:  &gen.Order{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.CreateOrder(ctx, tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(tt.expectedResp.Status, resp.Status)
				s.Equal(tt.expectedResp.TransactionId, resp.TransactionId)
			}
		})
	}
}

func (s *OrderServiceTestSuite) TestRemoveExpiryOrder() {
	userID := uuid.New()

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           time.Duration
		setupMock     func()
		expectedError string
		expectedResp  int
	}{
		{
			name: "Success",
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetExpiryOrders(gomock.Any(), 1*time.Minute).
					Return([]entity.Order{
						{
							ID:     uuid.New(),
							UserID: userID,
							Status: constanta.OrderStatusPending,
						},
						{
							ID:     uuid.New(),
							UserID: userID,
							Status: constanta.OrderStatusStockReserved,
						},
					}, nil)

				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(1), gomock.Any()).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusFailed, payloads["status"])
						return nil
					})

				s.mockWarehouseClient.EXPECT().
					ReleaseStock(gomock.Any(), gomock.Any()).
					Return(&gen.ReleaseStockResponse{
						ReleasedStockIds: []int64{1},
					}, nil)

				s.mockOrderRepo.EXPECT().
					UpdateOrderStatusWithCallback(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, status constanta.OrderStatus, id uuid.UUID, callback func() error) error {
					return callback()
				})

			},
			expectedError: "",
			expectedResp:  2,
			req:           1 * time.Minute,
		},
		{
			name: "Error when releasing stock",
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetExpiryOrders(gomock.Any(), 1*time.Minute).
					Return([]entity.Order{
						{
							ID:     uuid.New(),
							UserID: userID,
							Status: constanta.OrderStatusPending,
						},
						{
							ID:     uuid.New(),
							UserID: userID,
							Status: constanta.OrderStatusStockReserved,
						},
					}, nil)

				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(1), gomock.Any()).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusFailed, payloads["status"])
						return nil
					})

				s.mockWarehouseClient.EXPECT().
					ReleaseStock(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("failed to release stock"))

				s.mockOrderRepo.EXPECT().
					UpdateOrderStatusWithCallback(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, status constanta.OrderStatus, id uuid.UUID, callback func() error) error {
					return callback()
				})

			},
			expectedError: "",
			expectedResp:  2,
			req:           1 * time.Minute,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.RemoveExpiryOrder(ctx, tt.req)
			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Zero(resp)
			} else {
				s.NoError(err)
				s.NotZero(resp)
				s.Equal(tt.expectedResp, resp)
			}
		})
	}
}

func (s *OrderServiceTestSuite) TestCallbackTransaction() {
	userID := uuid.New()
	orderID := uuid.New()
	transactionID := "txn-abc-123"

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.CallbackTransactionRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Unknown status",
			req: &gen.CallbackTransactionRequest{
				TransactionId: transactionID,
				PaymentStatus: "UNKNOWN",
			},
			setupMock:     func() {},
			expectedError: "payment_status is UNKNOWN, must be one of PAID or FAILED",
		},
		{
			name: "Success with status Paid",
			req: &gen.CallbackTransactionRequest{
				TransactionId: transactionID,
				PaymentStatus: "PAID",
			},
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetOrderByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Order{
						ID:     orderID,
						Status: constanta.OrderStatusStockReserved,
					}, nil)

				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(1), orderID).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusCompleted, payloads["status"])
						return nil
					})

			},
			expectedError: "",
		},
		{
			name: "Success with status Failed",
			req: &gen.CallbackTransactionRequest{
				TransactionId: transactionID,
				PaymentStatus: "FAILED",
			},
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetOrderByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Order{
						ID:     orderID,
						Status: constanta.OrderStatusStockReserved,
					}, nil)

				s.mockOrderRepo.EXPECT().
					UpdateOrder(gomock.Any(), gomock.Len(1), orderID).
					DoAndReturn(func(ctx context.Context, payloads map[string]any, id uuid.UUID) error {
						s.Equal(constanta.OrderStatusFailed, payloads["status"])
						return nil
					})

			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.CallbackTransaction(ctx, tt.req)

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

func (s *OrderServiceTestSuite) TestGetOrder() {
	userID := uuid.New()
	orderID := uuid.New()

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.GetOrderRequest
		setupMock     func()
		expectedError string
	}{

		{
			name: "Success",
			req: &gen.GetOrderRequest{
				Id: orderID.String(),
			},
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Return(&entity.Order{
						ID:     orderID,
						UserID: userID,
					}, nil)

			},
			expectedError: "",
		},
		{
			name: "Permission denied",
			req: &gen.GetOrderRequest{
				Id: orderID.String(),
			},
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Return(&entity.Order{
						ID:     orderID,
						UserID: uuid.New(),
					}, nil)

			},
			expectedError: "you are not authorized to access this order",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetOrder(ctx, tt.req)

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

func (s *OrderServiceTestSuite) TestGetOrderList() {
	userID := uuid.New()

	md := metadata.New(map[string]string{
		string(globalcontanta.UserIDKey): userID.String(),
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	tests := []struct {
		name          string
		req           *gen.GetOrderListRequest
		setupMock     func()
		expectedError string
		expectedResp  *gen.Orders
	}{

		{
			name: "Success",
			req: &gen.GetOrderListRequest{
				StartDate: "2006-01-02",
				EndDate:   "2006-01-02",
				Status:    "PENDING",
			},
			setupMock: func() {
				s.mockOrderRepo.EXPECT().
					GetOrderList(gomock.Any(), gomock.Any()).
					Return([]entity.Order{
						{
							IdempotencyKey: uuid.UUID{},
							ID:             userID,
							UserID:         userID,
							Status:         "",
							TotalAmount:    &gen.Money{},
							TransactionID:  "",
							CreatedAt:      &time.Time{},
							UpdatedAt:      &time.Time{},
							Items:          []entity.OrderItem{},
						},
					}, nil)

			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetOrderList(ctx, tt.req)

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
