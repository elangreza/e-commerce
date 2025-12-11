package main

import (
	"context"
	"fmt"
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"log"

	"github.com/elangreza/e-commerce/order/internal/client"
	"github.com/elangreza/e-commerce/order/internal/server"
	"github.com/elangreza/e-commerce/order/internal/service"
	"github.com/elangreza/e-commerce/order/internal/sqlitedb"
	"github.com/elangreza/e-commerce/order/internal/task"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ServicePort          string `koanf:"SERVICE_PORT"`
	ProductServiceAddr   string `koanf:"PRODUCT_SERVICE_ADDR"`
	WarehouseServiceAddr string `koanf:"WAREHOUSE_SERVICE_ADDR"`
	ShopServiceAddr      string `koanf:"SHOP_SERVICE_ADDR"`
}

func main() {
	var cfg Config
	err := config.LoadConfig(&cfg)
	errChecker(err)

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

	// TODO payment service later
	paymentClient, err := client.NewPaymentClient()
	errChecker(err)

	// grpc clients
	grpcClientProduct, err := grpc.NewClient(cfg.ProductServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)
	grpcClientWarehouse, err := grpc.NewClient(cfg.WarehouseServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	orderService := service.NewOrderService(
		orderRepo,
		cartRepo,
		gen.NewWarehouseServiceClient(grpcClientWarehouse),
		gen.NewProductServiceClient(grpcClientProduct),
		paymentClient)

	srv := server.New(orderService)
	addr := fmt.Sprintf(":%s", cfg.ServicePort)
	go func() {
		if err := srv.Start(addr); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	fmt.Printf("ORDER-service running at %s\n", addr)

	taskOrder := task.NewTaskOrder(orderService)
	taskOrder.SetRemoveExpiryDuration(3 * time.Minute)

	gs := gracefulshutdown.New(context.Background(), 5*time.Second,
		gracefulshutdown.Operation{
			Name: "grpc",
			ShutdownFunc: func(ctx context.Context) error {
				srv.Close()
				return nil
			},
		},
		gracefulshutdown.Operation{
			Name: "task order",
			ShutdownFunc: func(ctx context.Context) error {
				taskOrder.Close()
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
