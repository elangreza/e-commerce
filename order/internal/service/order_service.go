package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	globalcontanta "github/elangreza/e-commerce/pkg/globalcontanta"

	"github.com/elangreza/e-commerce/gen"
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

	orderRepo interface{}

	stockServiceClient interface {
		GetStocks(ctx context.Context, productIds []string) (*gen.StockList, error)
		ReserveStock(ctx context.Context, cartItem []entity.CartItem) (*gen.ReserveStockResponse, error)
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
		return nil, err
	}

	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	reserveIDs, err := s.stockServiceClient.ReserveStock(ctx, cart.Items)
	if err != nil {
		return nil, errors.New("errors when reserving stocks")
	}

	fmt.Println(reserveIDs)

	// create order

	return nil, nil
}
