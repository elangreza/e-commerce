package service

//go:generate mockgen -source=product_service.go -destination=./mock/mock_product_service.go -package=mock

import (
	"context"
	"errors"

	"github.com/elangreza/e-commerce/pkg/globalcontanta"

	"github.com/elangreza/e-commerce/api/internal/constanta"
	params "github.com/elangreza/e-commerce/api/internal/params"
	"github.com/elangreza/e-commerce/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func NewOrderService(pClient gen.OrderServiceClient) *orderService {
	return &orderService{
		orderServiceClient: pClient,
	}
}

type orderService struct {
	orderServiceClient gen.OrderServiceClient
}

func (s *orderService) AddProductToCart(ctx context.Context, req params.AddToCartRequest) error {

	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return errors.New("error when parsing userID")
	}

	md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := s.orderServiceClient.AddProductToCart(newCtx, &gen.AddCartItemRequest{
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		return convertErrGrpc(err)
	}

	return nil
}

func (s *orderService) GetCart(ctx context.Context) (*params.GetCartResponse, error) {

	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return nil, errors.New("error when parsing userID")
	}

	md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	cart, err := s.orderServiceClient.GetCart(newCtx, &gen.Empty{})
	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.GetCartResponse{
		CartID: cart.Id,
		Items:  []params.GetCartItemsResponse{},
	}

	for _, item := range cart.Items {
		res.Items = append(res.Items, params.GetCartItemsResponse{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return res, nil
}

func (s *orderService) CreateOrder(ctx context.Context, req params.CreateOrderRequest) (*params.CreateOrderResponse, error) {

	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return nil, errors.New("error when parsing userID")
	}

	md := metadata.New(map[string]string{string(globalcontanta.UserIDKey): userID.String()})
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	order, err := s.orderServiceClient.CreateOrder(newCtx, &gen.CreateOrderRequest{
		IdempotencyKey: req.IdempotencyKey,
	})

	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.CreateOrderResponse{
		OrderID: order.Id,
		Items:   []params.GetCartItemsResponse{},
		TotalAmount: &params.Money{
			Units:        order.TotalAmount.Units,
			CurrencyCode: order.TotalAmount.CurrencyCode,
		},
		Status: order.Status,
	}

	for _, item := range order.Items {
		res.Items = append(res.Items, params.GetCartItemsResponse{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return res, nil
}
