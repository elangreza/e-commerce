package main

import (
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"
	"log"

	"github.com/elangreza/e-commerce/product/internal/client"
	"github.com/elangreza/e-commerce/product/internal/server"
	"github.com/elangreza/e-commerce/product/internal/service"
	"github.com/elangreza/e-commerce/product/internal/sqlitedb"
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

	address := fmt.Sprintf("localhost:%v", 50051)

	productService := service.NewProductService(productRepo, stockClient)
	srv := server.New(productService)
	if err := srv.Start(address); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
