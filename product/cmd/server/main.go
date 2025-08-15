package main

import (
	"fmt"
	"log"
	"net"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/grpcserver"
	"github.com/elangreza/e-commerce/product/internal/mockjson"
	"github.com/elangreza/e-commerce/product/internal/service"
	"google.golang.org/grpc"
)

func main() {

	// implement this later
	// github.com/samber/slog-zap

	productRepo, err := mockjson.LoadProductJson()
	errChecker(err)

	productService := service.NewProductService(productRepo)
	productServer := grpcserver.NewProductServer(productService)

	address := fmt.Sprintf(":%v", 50051)
	listener, err := net.Listen("tcp", address)
	errChecker(err)

	grpcServer := grpc.NewServer()
	gen.RegisterProductServiceServer(grpcServer, productServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}

func errChecker(err error) {
	if err != nil {
		panic(err)
	}
}
