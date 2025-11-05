package main

import (
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"
	"log"
	"net"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/product/internal/client"
	"github.com/elangreza/e-commerce/product/internal/grpcserver"
	"github.com/elangreza/e-commerce/product/internal/service"
	"github.com/elangreza/e-commerce/product/internal/sqlitedb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("product.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	// productRepo, err := mockjson.LoadProductJson()
	productRepo := sqlitedb.NewProductRepository(db)
	stockClient, err := client.NewStockClient()
	errChecker(err)

	productService := service.NewProductService(productRepo, stockClient)
	productServer := grpcserver.NewProductServer(productService)

	address := fmt.Sprintf("localhost:%v", 50051)
	listener, err := net.Listen("tcp", address)
	errChecker(err)

	grpcServer := grpc.NewServer()

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	gen.RegisterProductServiceServer(grpcServer, productServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
