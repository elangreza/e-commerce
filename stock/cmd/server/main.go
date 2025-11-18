package main

import (
	"github/elangreza/e-commerce/pkg/dbsql"
	"github/elangreza/e-commerce/stock/internal/server"
	"github/elangreza/e-commerce/stock/internal/service"
	"github/elangreza/e-commerce/stock/internal/sqlitedb"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	address := "localhost:50052"

	srv := server.New(stockService)

	if err := srv.Start(address); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}

	// grpcServer := grpc.NewServer(
	// 	grpc.ChainUnaryInterceptor(
	// 		interceptor.UserIDParser(),
	// 	),
	// )

}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
