package main

import (
	"fmt"
	"log"

	handler "github.com/elangreza/e-commerce/api/internal/handlers"
	"github.com/elangreza/e-commerce/api/internal/server"
	"github.com/elangreza/e-commerce/api/internal/service"
	"github.com/elangreza/e-commerce/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	grpcClient, err := grpc.NewClient("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	productService := service.NewProductService(gen.NewProductServiceClient(grpcClient))
	handler := handler.NewProductHandler(productService)
	srv := server.New(":8080", handler)
	if err := srv.Start(); err != nil {
		fmt.Printf("err %v\n", err)
		return
	}
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
