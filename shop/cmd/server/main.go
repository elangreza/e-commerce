package main

import (
	"context"
	"log"
	"time"

	"github.com/elangreza/e-commerce/shop/internal/client"
	"github.com/elangreza/e-commerce/shop/internal/server"
	"github.com/elangreza/e-commerce/shop/internal/service"
	"github.com/elangreza/e-commerce/shop/internal/sqlitedb"

	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("shop.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
		dbsql.WithAutoSeeder("file://./migrations/seed"),
	)
	errChecker(err)
	defer db.Close()

	warehouseClient, err := client.NewWarehouseClient()
	errChecker(err)

	shopRepo := sqlitedb.NewShopRepo(db)
	shopService := service.NewShopService(shopRepo, warehouseClient)

	address := "localhost:50054"

	srv := server.New(shopService)
	go func() {
		if err := srv.Start(address); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

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
