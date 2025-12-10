package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elangreza/e-commerce/shop/internal/client"
	"github.com/elangreza/e-commerce/shop/internal/server"
	"github.com/elangreza/e-commerce/shop/internal/service"
	"github.com/elangreza/e-commerce/shop/internal/sqlitedb"

	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"

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

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("shop.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
		dbsql.WithAutoSeeder("file://./migrations/seed"),
	)
	errChecker(err)
	defer db.Close()

	warehouseClient, err := client.NewWarehouseClient(cfg.WarehouseServiceAddr)
	errChecker(err)

	shopRepo := sqlitedb.NewShopRepo(db)
	shopService := service.NewShopService(shopRepo, warehouseClient)

	addr := fmt.Sprintf(":%s", cfg.ServicePort)

	srv := server.New(shopService)
	go func() {
		if err := srv.Start(addr); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	fmt.Printf("SHOP-service running at %s\n", addr)

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
