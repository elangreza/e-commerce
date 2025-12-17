package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/elangreza/e-commerce/pkg/contextrequest"
	"github.com/elangreza/e-commerce/pkg/extractor"
	"github.com/elangreza/e-commerce/pkg/money"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/constanta"
	"github.com/elangreza/e-commerce/order/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=order_service.go -destination=mock/mock_order_service.go -package=mock
//go:generate mockgen -package=mock -destination=mock/mock_deps.go github.com/elangreza/e-commerce/gen WarehouseServiceClient,PaymentServiceClient,ProductServiceClient

type (
	cartRepo interface {
		GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
		CreateCart(ctx context.Context, cart entity.Cart) error
		UpdateCartItem(ctx context.Context, item entity.CartItem) error
	}

	orderRepo interface {
		CreateOrder(ctx context.Context, order entity.Order) (uuid.UUID, error)
		GetOrderByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*entity.Order, error)
		UpdateOrder(ctx context.Context, payloads map[string]any, orderID uuid.UUID) error
		GetExpiryOrders(ctx context.Context, duration time.Duration) ([]entity.Order, error)
		UpdateOrderStatusWithCallback(ctx context.Context, status constanta.OrderStatus, orderID uuid.UUID, callback func() error) error
		GetOrderByTransactionID(ctx context.Context, transactionID string) (*entity.Order, error)
		GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)
		GetOrderList(ctx context.Context, req entity.GetOrderListRequest) ([]entity.Order, error)
	}
)

type OrderService struct {
	orderRepo              orderRepo
	cartRepo               cartRepo
	warehouseServiceClient gen.WarehouseServiceClient
	productServiceClient   gen.ProductServiceClient
	paymentServiceClient   gen.PaymentServiceClient
	gen.UnimplementedOrderServiceServer
}

func NewOrderService(
	orderRepo orderRepo,
	cartRepo cartRepo,
	warehouseServiceClient gen.WarehouseServiceClient,
	productServiceClient gen.ProductServiceClient,
	paymentServiceClient gen.PaymentServiceClient,
) *OrderService {
	return &OrderService{
		orderRepo:              orderRepo,
		cartRepo:               cartRepo,
		warehouseServiceClient: warehouseServiceClient,
		productServiceClient:   productServiceClient,
		paymentServiceClient:   paymentServiceClient,
	}
}

func (s *OrderService) AddProductToCart(ctx context.Context, req *gen.AddCartItemRequest) (*gen.Empty, error) {
	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	withStock := true
	products, err := s.productServiceClient.GetProducts(ctx, &gen.GetProductsRequest{
		Ids:       []string{req.ProductId},
		WithStock: withStock,
	})
	if err != nil {
		return nil, err
	}
	if products == nil || products.Products == nil || len(products.Products) == 0 {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	product := products.Products[0]

	if req.Quantity > product.Stock {
		return nil, status.Errorf(codes.InvalidArgument, "quantity cannot exceed the maximum stock, current stock is %d", product.Stock)
	}

	if cart == nil {
		cart = &entity.Cart{
			UserID: userID,
			Items: []entity.CartItem{
				{
					ProductID: req.ProductId,
					Quantity:  req.Quantity,
					Name:      product.GetName(),
					Price:     product.GetPrice(),
				},
			},
		}

		err = s.cartRepo.CreateCart(ctx, *cart)
		if err != nil {
			return nil, err
		}

		// return early after cart creation
		return &gen.Empty{}, nil
	}

	err = s.cartRepo.UpdateCartItem(ctx, entity.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductId,
		Quantity:  req.Quantity,
		Name:      product.GetName(),
		Price:     product.GetPrice(),
	})
	if err != nil {
		return nil, err
	}

	return &gen.Empty{}, nil
}

func (s *OrderService) GetCart(ctx context.Context, req *gen.Empty) (*gen.Cart, error) {
	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "cart not found")
		}
		return nil, err
	}

	productIDs := make([]string, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	stocks, err := s.warehouseServiceClient.GetStocks(ctx, &gen.GetStockRequest{
		ProductIds: productIDs,
	})
	if err != nil {
		return nil, err
	}

	stockMap := make(map[string]int64)
	for _, stock := range stocks.Stocks {
		stockMap[stock.ProductId] = stock.Quantity
	}

	for i, item := range cart.Items {
		if stock, ok := stockMap[item.ProductID]; ok {
			cart.Items[i].ActualStock = stock
		} else {
			cart.Items[i].ActualStock = 0
		}
	}

	return cart.GetGenCart(), nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req *gen.CreateOrderRequest) (*gen.Order, error) {
	idempotencyKey, err := uuid.Parse(req.IdempotencyKey)
	if err != nil {
		return nil, errors.New("invalid idempotency_key format")
	}

	ord, err := s.orderRepo.GetOrderByIdempotencyKey(ctx, idempotencyKey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if ord != nil {
		return ord.GetGenOrder(), nil
	}

	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, status.Errorf(codes.NotFound, "cart not found")
	}

	for _, item := range cart.Items {
		if item.Quantity <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "Order item quantity must be greater than 0")
		}
	}

	// Enforce single-currency cart (required to safely sum totalAmount)
	var cartCurrency string
	orderItems := make([]entity.OrderItem, 0, len(cart.Items))
	totalAmount, _ := money.New(0, "IDR")

	var withStock = false
	products, err := s.productServiceClient.GetProducts(ctx, &gen.GetProductsRequest{
		Ids:       cart.GetProductIDs(),
		WithStock: withStock,
	})
	if err != nil {
		return nil, errors.New("failed to fetch products")
	}

	productsMap := make(map[string]*gen.Product)
	for _, product := range products.Products {
		productsMap[product.Id] = product
	}

	for _, item := range cart.Items {
		product, ok := productsMap[item.ProductID]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}

		price := product.GetPrice()
		if price == nil {
			return nil, fmt.Errorf("product %s has no price", item.ProductID)
		}

		// Validate currency consistency
		if cartCurrency == "" {
			cartCurrency = price.GetCurrencyCode()
		} else if cartCurrency != price.GetCurrencyCode() {
			return nil, status.Errorf(codes.InvalidArgument, "mixed currencies in cart are not supported")
		}

		totalPricePerUnit, err := money.MultiplyByInt(price, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total price for product %s: %w", item.ProductID, err)
		}

		orderItems = append(orderItems, entity.OrderItem{
			ProductID:         item.ProductID,
			Name:              product.GetName(),
			PricePerUnit:      price,
			Quantity:          item.Quantity,
			TotalPricePerUnit: totalPricePerUnit,
		})

		totalAmount, err = money.Add(totalAmount, totalPricePerUnit)
		if err != nil {
			return nil, fmt.Errorf("failed to accumulate total amount: %w", err)
		}
	}

	order := entity.Order{
		IdempotencyKey: idempotencyKey,
		UserID:         userID,
		Status:         constanta.OrderStatusPending, // New initial status
		Items:          orderItems,
		TotalAmount:    totalAmount,
	}

	orderID, err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to persist order: %w", err)
	}

	// Create a rollback function for cleanup
	rollback := func() error {
		return s.orderRepo.UpdateOrder(ctx, map[string]any{
			"status": constanta.OrderStatusFailed,
		}, orderID)
	}

	ctx = contextrequest.AppendUserIDintoContextGrpcClient(ctx, userID)

	stocks := []*gen.Stock{}
	for _, item := range cart.Items {
		stocks = append(stocks, &gen.Stock{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	// Reserve stock
	_, err = s.warehouseServiceClient.ReserveStock(ctx, &gen.ReserveStockRequest{
		OrderId: orderID.String(),
		Stocks:  stocks,
	})
	if err != nil {
		rollbackErr := rollback()
		if rollbackErr != nil {
			// Log this error - partial failure state
			fmt.Printf("Error during rollback: %v", rollbackErr)
		}
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Process payment
	paymentTransaction, err := s.paymentServiceClient.ProcessPayment(ctx, &gen.ProcessPaymentRequest{
		OrderId:     orderID.String(),
		TotalAmount: totalAmount,
	})
	if err != nil {
		// Release reserved stock
		_, releaseErr := s.warehouseServiceClient.ReleaseStock(ctx, &gen.ReleaseStockRequest{
			OrderId: orderID.String(),
		})
		if releaseErr != nil {
			// Log - stock might be stuck in reserved state
			fmt.Printf("Error releasing stock during payment failure: %v", releaseErr)
		}

		rollbackErr := rollback()
		if rollbackErr != nil {
			fmt.Printf("Error during rollback: %v", rollbackErr)
		}

		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Update order status to indicate stock is reserved and include transaction ID
	err = s.orderRepo.UpdateOrder(ctx, map[string]any{
		"status":         constanta.OrderStatusStockReserved,
		"transaction_id": paymentTransaction.TransactionId,
	}, orderID)
	if err != nil {
		// Payment succeeded but status update failed - log this inconsistency
		fmt.Printf("Payment succeeded but status update failed: %v", err)
		// Consider whether to return an error or continue with partial success
		// For now, we'll return the error to indicate the operation didn't complete successfully
		return nil, fmt.Errorf("failed to update order with transaction ID: %w", err)
	}

	order.ID = orderID
	order.Status = constanta.OrderStatusStockReserved
	order.TransactionID = paymentTransaction.TransactionId

	return order.GetGenOrder(), nil
}

func (s *OrderService) RemoveExpiryOrder(ctx context.Context, duration time.Duration) (int, error) {
	orders, err := s.orderRepo.GetExpiryOrders(ctx, duration)
	if err != nil {
		return 0, err
	}

	for _, order := range orders {

		if order.Status == constanta.OrderStatusPending {
			// TODO must be retry flow, process the cause of pending order
			// if the cause is stock reservation failed, then release stock
			// if the cause is payment failed, then update status to failed
			// if the cause is repository is failed, then update status to failed
			// if the cause is unknown, then update status to failed
			err = s.orderRepo.UpdateOrder(ctx, map[string]any{
				"status": constanta.OrderStatusFailed,
			}, order.ID)
			if err != nil {
				fmt.Println("err when Update status", err)
			}
		}

		if order.Status == constanta.OrderStatusStockReserved {
			var releaseStock = func() error {
				ctx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), order.UserID)
				_, err := s.warehouseServiceClient.ReleaseStock(ctx, &gen.ReleaseStockRequest{
					OrderId: order.ID.String(),
				})
				if err != nil {
					return err
				}

				return nil
			}

			err := s.orderRepo.UpdateOrderStatusWithCallback(ctx, constanta.OrderStatusFailed, order.ID, releaseStock)
			if err != nil {
				fmt.Println("err when Release Stock", err)
			}

		}
	}

	return len(orders), nil
}

func (s *OrderService) CallbackTransaction(ctx context.Context, req *gen.CallbackTransactionRequest) (*gen.Empty, error) {
	if req.PaymentStatus == "" {
		return nil, status.Errorf(codes.InvalidArgument, "payment_status cannot be empty")
	}

	var paymentStatus constanta.PaymentStatus
	err := paymentStatus.Scan(req.PaymentStatus)
	if err != nil {
		return nil, err
	}

	if paymentStatus.String() == "UNKNOWN" {
		return nil, status.Errorf(codes.InvalidArgument, "payment_status is %s, must be one of %s or %s", paymentStatus.String(), constanta.PAID, constanta.FAILED)
	}

	order, err := s.orderRepo.GetOrderByTransactionID(ctx, req.TransactionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	if order.Status != constanta.OrderStatusStockReserved {
		return nil, status.Errorf(codes.FailedPrecondition, "order status must be status reserved")
	}

	if paymentStatus == constanta.PAID {
		err = s.orderRepo.UpdateOrder(ctx, map[string]any{
			"status": constanta.OrderStatusCompleted,
		}, order.ID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update order: %v", err)
		}
	}

	if paymentStatus == constanta.FAILED {
		err = s.orderRepo.UpdateOrder(ctx, map[string]any{
			"status": constanta.OrderStatusFailed,
		}, order.ID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update order: %v", err)
		}
	}

	return &gen.Empty{}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *gen.GetOrderRequest) (*gen.Order, error) {
	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order id")
	}

	order, err := s.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	if order.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not authorized to access this order")
	}

	return order.GetGenOrder(), nil
}

func (s *OrderService) GetOrderList(ctx context.Context, req *gen.GetOrderListRequest) (*gen.Orders, error) {
	userID, err := extractor.ExtractUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	var isFilterByStatus bool
	var reqStatus constanta.OrderStatus
	if req.Status != "" {
		err = reqStatus.Scan(req.Status)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "status not valid")
		}

		if reqStatus.String() == "UNKNOWN" {
			return nil, status.Errorf(codes.InvalidArgument, "status not valid")
		}

		isFilterByStatus = true
	}

	var isFilterByDate bool
	var startDate time.Time
	var endDate time.Time
	if req.StartDate != "" && req.EndDate != "" {
		startDate, err = time.Parse(time.DateOnly, req.StartDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "start_date not valid")
		}

		endDate, err = time.Parse(time.DateOnly, req.EndDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "end_date not valid")
		}

		if startDate.After(endDate) {
			return nil, status.Errorf(codes.InvalidArgument, "start_date must be before end_date")
		}

		isFilterByDate = true
	}

	orderList, err := s.orderRepo.GetOrderList(ctx, entity.GetOrderListRequest{
		UserID:           userID,
		IsFilterByDate:   isFilterByDate,
		StartDate:        startDate,
		EndDate:          endDate,
		IsFilterByStatus: isFilterByStatus,
		Status:           reqStatus,
	})
	if err != nil {
		return nil, err
	}

	orders := []*gen.Order{}
	for _, order := range orderList {
		orders = append(orders, order.GetGenOrder())
	}

	return &gen.Orders{
		Orders: orders,
	}, nil
}
