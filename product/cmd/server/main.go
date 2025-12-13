package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
		dbsql.WithSqliteDB("data/product.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
		dbsql.WithAutoSeeder("file://./migrations/seed"),
	)
	errChecker(err)

	productRepo := sqlitedb.NewProductRepository(db)

	// warehouse
	grpcClientWarehouse, err := grpc.NewClient(cfg.WarehouseServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	productService := service.NewProductService(productRepo, gen.NewWarehouseServiceClient(grpcClientWarehouse))

	addr := fmt.Sprintf(":%s", cfg.ServicePort)
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
