package grpcserver

// go generate
//go:generate mockgen -source=order_grpc.go -destination=./mock/mock_order_grpc.go -package=mock

import (
	"context"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	orderService interface {
		AddProductToCart(ctx context.Context, productId string, quantity int64) error
		GetCart(ctx context.Context) (*entity.Cart, error)
		CreateOrder(ctx context.Context, idempotencyKey string) (*entity.Order, error)
	}

	OrderServer struct {
		orderService orderService
		gen.UnimplementedOrderServiceServer
	}
)

func NewOrderServer(orderService orderService) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}

func (o *OrderServer) AddProductToCart(ctx context.Context, req *gen.AddCartItemRequest) (*gen.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProductToCart not implemented")
}

func (o *OrderServer) GetCart(ctx context.Context, req *gen.Empty) (*gen.Cart, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCart not implemented")
}

func (o *OrderServer) CreateOrder(ctx context.Context, req *gen.CreateOrderRequest) (*gen.Order, error) {

	order, err := o.orderService.CreateOrder(ctx, req.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	items := []*gen.OrderItem{}
	for _, item := range order.Items {
		items = append(items, &gen.OrderItem{
			ProductId:    item.ProductID,
			Name:         item.Name,
			PricePerUnit: item.PricePerUnit,
			Quantity:     item.Quantity,
		})
	}

	return &gen.Order{
		Id:          order.ID.String(),
		UserId:      order.UserID.String(),
		Items:       items,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
	}, nil
}
