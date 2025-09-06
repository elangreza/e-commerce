package main

import (
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"
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

	orderRepo := sqlitedb.NewOrderRepository(db)
	stockClient, err := client.NewStockClient()
	errChecker(err)

	orderService := service.NewOrderService(orderRepo, stockClient)
	orderServer := grpcserver.NewOrderServer(orderService)

	address := fmt.Sprintf(":%v", 50051)
	listener, err := net.Listen("tcp", address)
	errChecker(err)

	grpcServer := grpc.NewServer()
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
