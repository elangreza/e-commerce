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
)

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

type (
	cartRepo interface {
		GetCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
		CreateCart(ctx context.Context, cart entity.Cart) error
		UpdateCartItem(ctx context.Context, cartID uuid.UUID, item entity.CartItem) error
	}

	orderRepo interface {
		CreateOrder(ctx context.Context, order entity.Order) (uuid.UUID, error)
	}

	stockServiceClient interface {
		GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error)
		ReserveStock(ctx context.Context, cartItem []entity.CartItem) (*gen.ReserveStockResponse, error)
		ReleaseStock(ctx context.Context, reservedStockIds []int64) (*gen.ReleaseStockResponse, error)
	}

	productServiceClient interface {
		GetProduct(ctx context.Context, productId string) (*gen.Product, error)
	}
)

type orderService struct {
	orderRepo            orderRepo
	cartRepo             cartRepo
	stockServiceClient   stockServiceClient
	productServiceClient productServiceClient
}

func NewOrderService(
	orderRepo orderRepo,
	cartRepo cartRepo,
	stockServiceClient stockServiceClient,
	productServiceClient productServiceClient) *orderService {
	return &orderService{
		orderRepo:            orderRepo,
		cartRepo:             cartRepo,
		stockServiceClient:   stockServiceClient,
		productServiceClient: productServiceClient,
	}
}

func (s *orderService) AddProductToCart(ctx context.Context, productId string, quantity int64) error {
	userID, ok := ctx.Value(globalcontanta.UserIDKey).(uuid.UUID)
	if !ok {
		return errors.New("unauthorized")
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if cart == nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		product, err := s.productServiceClient.GetProduct(ctx, productId)
		if err != nil {
			return err
		}
		if product == nil {
			return errors.New("product not found")
		}

		cart = &entity.Cart{
			ID:     id,
			UserID: userID,
			Items: []entity.CartItem{
				{
					ProductID: productId,
					Quantity:  quantity,
					Price:     product.Price,
				},
			},
		}

		err = s.cartRepo.CreateCart(ctx, *cart)
		if err != nil {
			return err
		}
		return nil
	}

	err = s.cartRepo.UpdateCartItem(ctx, cart.ID, entity.CartItem{
		ProductID: productId,
		Quantity:  quantity,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *orderService) GetCart(ctx context.Context) (*entity.Cart, error) {
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
		return &entity.Cart{
			Items: []entity.CartItem{},
		}, nil
	}

	if len(cart.Items) == 0 {
		return &entity.Cart{
			Items: []entity.CartItem{},
		}, nil
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
			cart.Items[i].Stock = stock
		} else {
			cart.Items[i].Stock = 0
		}
	}

	return cart, nil
}

func (s *orderService) CreateOrder(ctx context.Context) (*entity.Order, error) {
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

	reserveIDs, err := s.stockServiceClient.ReserveStock(ctx, cart.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Defer stock release on any error after reservation
	var finalOrder *entity.Order
	var finalErr error
	defer func() {
		if finalErr != nil {
			_, releaseErr := s.stockServiceClient.ReleaseStock(ctx, reserveIDs.GetReservedStockIds())
			if releaseErr != nil {
				// TODO: log error (e.g., s.logger.Error("failed to release reserved stock", ...))
			}
		}
	}()

	// Enforce single-currency cart (required to safely sum totalAmount)
	var cartCurrency string
	orderItems := make([]entity.OrderItem, 0, len(cart.Items))
	totalAmount := &gen.Money{}

	for _, item := range cart.Items {
		product, err := s.productServiceClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			finalErr = fmt.Errorf("failed to fetch product %s: %w", item.ProductID, err)
			return nil, finalErr
		}

		price := product.GetPrice()
		if price == nil {
			finalErr = fmt.Errorf("product %s has no price", item.ProductID)
			return nil, finalErr
		}

		// Validate currency consistency
		if cartCurrency == "" {
			cartCurrency = price.GetCurrencyCode()
		} else if cartCurrency != price.GetCurrencyCode() {
			finalErr = errors.New("mixed currencies in cart are not supported")
			return nil, finalErr
		}

		totalPricePerUnit, err := money.MultiplyByInt(price, item.Quantity)
		if err != nil {
			finalErr = fmt.Errorf("failed to calculate total price for product %s: %w", item.ProductID, err)
			return nil, finalErr
		}

		orderItems = append(orderItems, entity.OrderItem{
			ProductID:         item.ProductID,
			Name:              product.GetName(),
			PricePerUnit:      price,
			Currency:          price.GetCurrencyCode(),
			Quantity:          item.Quantity,
			TotalPricePerUnit: totalPricePerUnit,
		})

		totalAmount, err = money.Add(totalAmount, totalPricePerUnit)
		if err != nil {
			finalErr = fmt.Errorf("failed to accumulate total amount: %w", err)
			return nil, finalErr
		}
	}

	order := entity.Order{
		UserID:      userID,
		Status:      constanta.OrderStatusPending,
		Items:       orderItems,
		TotalAmount: totalAmount,
		Currency:    cartCurrency,
	}

	orderID, err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		finalErr = fmt.Errorf("failed to persist order: %w", err)
		return nil, finalErr
	}

	order.ID = orderID
	finalOrder = &order
	return finalOrder, nil
}
