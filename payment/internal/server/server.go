package server

import (
	"net"

	"github.com/elangreza/e-commerce/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	service    gen.PaymentServiceServer
}

func New(svc gen.PaymentServiceServer) *Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	gen.RegisterPaymentServiceServer(grpcServer, svc)

	return &Server{
		grpcServer: grpcServer,
		service:    svc,
	}
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(lis)
}
