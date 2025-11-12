package main

import (
	"github/elangreza/e-commerce/pkg/dbsql"
	"github/elangreza/e-commerce/pkg/interceptor"
	"github/elangreza/e-commerce/stock/internal/grpcserver"
	"github/elangreza/e-commerce/stock/internal/service"
	"github/elangreza/e-commerce/stock/internal/sqlitedb"
	"log"
	"net"

	"github.com/elangreza/e-commerce/gen"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("stock.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	stockRepo := sqlitedb.NewStockRepo(db)
	stockService := service.NewStockService(stockRepo)
	stockServer := grpcserver.NewStockGRPCServer(stockService)

	address := "localhost:50052"
	listener, err := net.Listen("tcp", address)
	errChecker(err)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UserIDParser(),
		),
	)

	gen.RegisterStockServiceServer(grpcServer, stockServer)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
