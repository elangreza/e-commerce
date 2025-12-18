package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elangreza/e-commerce/warehouse/internal/server"
	"github.com/elangreza/e-commerce/warehouse/internal/service"
	"github.com/elangreza/e-commerce/warehouse/internal/sqlitedb"

	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ServicePort string `koanf:"SERVICE_PORT"`
	DBPath      string `koanf:"DB_PATH"`
}

func main() {
	var cfg Config
	err := config.LoadConfig(&cfg)
	errChecker(err)

	// Default to data-local for local development
	dbPath := cfg.DBPath
	if dbPath == "" {
		dbPath = "data-local/warehouse.db"
	}

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB(dbPath),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
		dbsql.WithAutoSeeder("file://./migrations/seed"),
	)
	errChecker(err)

	warehouseRepo := sqlitedb.NewWarehouseRepo(db)
	warehouseService := service.NewWarehouseService(warehouseRepo)

	addr := fmt.Sprintf(":%s", cfg.ServicePort)

	srv := server.New(warehouseService)
	go func() {
		if err := srv.Start(addr); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	fmt.Printf("WAREHOUSE-service running at %s\n", addr)

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
