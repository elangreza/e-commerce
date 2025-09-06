package grpcserver

// go generate
//go:generate mockgen -source=order_grpc.go -destination=./mock/mock_order_grpc.go -package=mock

import (
	"github.com/elangreza/e-commerce/gen"
)

type (
	orderService interface {
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
