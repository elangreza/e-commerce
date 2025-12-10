package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"

	"github.com/elangreza/e-commerce/product/internal/client"
	"github.com/elangreza/e-commerce/product/internal/server"
	"github.com/elangreza/e-commerce/product/internal/service"
	"github.com/elangreza/e-commerce/product/internal/sqlitedb"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ServicePort          string `koanf:"SERVICE_PORT"`
	WarehouseServiceAddr string `koanf:"WAREHOUSE_SERVICE_ADDR"`
}

func main() {
	var cfg Config
	err := config.LoadConfig(&cfg)
	errChecker(err)

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("product.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
		dbsql.WithAutoSeeder("file://./migrations/seed"),
	)
	errChecker(err)
	defer db.Close()

	productRepo := sqlitedb.NewProductRepository(db)
	warehouseClient, err := client.NewWarehouseClient(cfg.WarehouseServiceAddr)
	errChecker(err)

	addr := fmt.Sprintf(":%s", cfg.ServicePort)

	productService := service.NewProductService(productRepo, warehouseClient)
	srv := server.New(productService)
	go func() {
		if err := srv.Start(addr); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	fmt.Printf("PRODUCT-service running at %s\n", addr)

	gs := gracefulshutdown.New(context.Background(), 5*time.Second,
		gracefulshutdown.Operation{
			Name: "grpc",
			ShutdownFunc: func(ctx context.Context) error {
				srv.Close()
				return nil
			},
		},
		gracefulshutdown.Operation{
			Name: "sqlite",
			ShutdownFunc: func(ctx context.Context) error {
				return db.Close()
			},
		},
	)
	<-gs
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
