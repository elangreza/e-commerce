package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	globalcontanta "github/elangreza/e-commerce/pkg/globalcontanta"
	"github/elangreza/e-commerce/pkg/money"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/constanta"
	"github.com/elangreza/e-commerce/order/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

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
	}

	stockServiceClient interface {
		GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error)
		ReserveStock(ctx context.Context, cartItem []entity.CartItem) (*gen.ReserveStockResponse, error)
		ReleaseStock(ctx context.Context, reservedStockIds []int64) (*gen.ReleaseStockResponse, error)
	}

	productServiceClient interface {
		GetProduct(ctx context.Context, productId string) (*gen.Product, error)
		GetProducts(ctx context.Context, withStock bool, productIds ...string) (*gen.Products, error)
	}

	paymentServiceClient interface {
		ProcessPayment(ctx context.Context, totalAmount *gen.Money, orderID uuid.UUID) (*gen.ProcessPaymentResponse, error)
	}
)

type orderService struct {
	orderRepo            orderRepo
	cartRepo             cartRepo
	stockServiceClient   stockServiceClient
	productServiceClient productServiceClient
	paymentServiceClient paymentServiceClient
	gen.UnimplementedOrderServiceServer
}

func NewOrderService(
	orderRepo orderRepo,
	cartRepo cartRepo,
	stockServiceClient stockServiceClient,
	productServiceClient productServiceClient,
	paymentServiceClient paymentServiceClient,
) *orderService {
	return &orderService{
		orderRepo:            orderRepo,
		cartRepo:             cartRepo,
		stockServiceClient:   stockServiceClient,
		productServiceClient: productServiceClient,
		paymentServiceClient: paymentServiceClient,
	}
}

func (s *orderService) AddProductToCart(ctx context.Context, req *gen.AddCartItemRequest) (*gen.Empty, error) {
	userID, ok := ctx.Value(globalcontanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	product, err := s.productServiceClient.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
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

func (s *orderService) GetCart(ctx context.Context, req *gen.Empty) (*gen.Cart, error) {
	userID, ok := ctx.Value(globalcontanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		err = s.cartRepo.CreateCart(ctx, entity.Cart{
			ID:     uuid.Must(uuid.NewV7()),
			UserID: userID,
			Items:  []entity.CartItem{},
		})
		if err != nil {
			return nil, err
		}
		return cart.GetGenCart(), nil
	}

	if len(cart.Items) == 0 {
		return cart.GetGenCart(), nil
	}

	productIDs := make([]string, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	stocks, err := s.stockServiceClient.GetStocks(ctx, productIDs)
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

func (s *orderService) CreateOrder(ctx context.Context, req *gen.CreateOrderRequest) (*gen.Order, error) {
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

	userID, ok := ctx.Value(globalcontanta.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
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
	products, err := s.productServiceClient.GetProducts(ctx, withStock, cart.GetProductIDs()...)
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
			return nil, fmt.Errorf("failed to fetch product %s", item.ProductID)
		}

		price := product.GetPrice()
		if price == nil {
			return nil, fmt.Errorf("product %s has no price", item.ProductID)
		}

		// Validate currency consistency
		if cartCurrency == "" {
			cartCurrency = price.GetCurrencyCode()
		} else if cartCurrency != price.GetCurrencyCode() {
			return nil, errors.New("mixed currencies in cart are not supported")
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

	// Reserve stock
	reserveIDs, err := s.stockServiceClient.ReserveStock(ctx, cart.Items)
	if err != nil {
		rollbackErr := rollback()
		if rollbackErr != nil {
			// Log this error - partial failure state
			fmt.Printf("Error during rollback: %v", rollbackErr)
		}
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Process payment
	paymentTransaction, err := s.paymentServiceClient.ProcessPayment(ctx, order.TotalAmount, orderID)
	if err != nil {
		// Release reserved stock
		_, releaseErr := s.stockServiceClient.ReleaseStock(ctx, reserveIDs.GetReservedStockIds())
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
		return nil, fmt.Errorf("payment succeeded but failed to update order with transaction ID: %w", err)
	}

	order.ID = orderID
	return ord.GetGenOrder(), nil
}
