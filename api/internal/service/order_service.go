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

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

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

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

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

func (s *orderService) CreateOrder(ctx context.Context, req params.CreateOrderRequest) (*params.OrderResponse, error) {

	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return nil, errors.New("error when parsing userID")
	}

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

	order, err := s.orderServiceClient.CreateOrder(newCtx, &gen.CreateOrderRequest{
		IdempotencyKey: req.IdempotencyKey,
	})

	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.OrderResponse{
		OrderID: order.GetId(),
		Items:   []params.GetCartItemsResponse{},
		TotalAmount: &params.Money{
			Units:        order.GetTotalAmount().GetUnits(),
			CurrencyCode: order.GetTotalAmount().GetCurrencyCode(),
		},
		Status:        order.GetStatus(),
		TransactionID: order.GetTransactionId(),
	}

	for _, item := range order.Items {
		res.Items = append(res.Items, params.GetCartItemsResponse{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		})
	}

	return res, nil
}

func (s *orderService) GetOrderList(ctx context.Context, req params.GetOrderListRequest) (*params.GetOrderListResponse, error) {
	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return nil, errors.New("error when parsing userID")
	}

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

	list, err := s.orderServiceClient.GetOrderList(newCtx, &gen.GetOrderListRequest{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    req.Status,
	})

	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.GetOrderListResponse{
		OrderList: []params.OrderResponse{},
	}

	for _, item := range list.GetOrders() {
		res.OrderList = append(res.OrderList, params.OrderResponse{
			OrderID: item.GetId(),
			TotalAmount: &params.Money{
				Units:        item.GetTotalAmount().GetUnits(),
				CurrencyCode: item.GetTotalAmount().GetCurrencyCode(),
			},
			Status:        item.GetStatus(),
			TransactionID: item.GetTransactionId(),
			Items:         nil,
		})
	}

	return res, nil
}

func (s *orderService) GetOrderDetail(ctx context.Context, orderID string) (*params.OrderResponse, error) {
	userID, ok := ctx.Value(constanta.LocalUserID).(uuid.UUID)
	if !ok {
		return nil, errors.New("error when parsing userID")
	}

	newCtx := contextrequest.AppendUserIDintoContextGrpcClient(context.Background(), userID)

	order, err := s.orderServiceClient.GetOrder(newCtx, &gen.GetOrderRequest{
		Id: orderID,
	})

	if err != nil {
		return nil, convertErrGrpc(err)
	}

	res := &params.OrderResponse{
		OrderID: order.GetId(),
		Items:   []params.GetCartItemsResponse{},
		TotalAmount: &params.Money{
			Units:        order.GetTotalAmount().GetUnits(),
			CurrencyCode: order.GetTotalAmount().GetCurrencyCode(),
		},
		Status:        order.GetStatus(),
		TransactionID: order.GetTransactionId(),
	}

	for _, item := range order.Items {
		res.Items = append(res.Items, params.GetCartItemsResponse{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		})
	}

	return res, nil
}
