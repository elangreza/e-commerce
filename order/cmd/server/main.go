package main

import (
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"

	"github/elangreza/e-commerce/pkg/interceptor"
	"log"
	"net"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/order/internal/client"
	"github.com/elangreza/e-commerce/order/internal/grpcserver"
	"github.com/elangreza/e-commerce/order/internal/service"
	"github.com/elangreza/e-commerce/order/internal/sqlitedb"
	"google.golang.org/grpc"
)

func main() {

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("order.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	cartRepo := sqlitedb.NewCartRepository(db)
	orderRepo := sqlitedb.NewOrderRepository(db)
	stockClient, err := client.NewStockClient()
	errChecker(err)
	productClient, err := client.NewProductClient()
	errChecker(err)
	paymentClient, err := client.NewPaymentClient()
	errChecker(err)

	orderService := service.NewOrderService(orderRepo, cartRepo, stockClient, productClient, paymentClient)
	orderServer := grpcserver.NewOrderServer(orderService)

	address := fmt.Sprintf(":%v", 50051)
	listener, err := net.Listen("tcp", address)
	errChecker(err)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UserIDParser(),
		),
	)
	gen.RegisterOrderServiceServer(grpcServer, orderServer)
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
