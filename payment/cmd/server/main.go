package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/elangreza/e-commerce/payment/internal/server"
	"github.com/elangreza/e-commerce/payment/internal/service"
	"github.com/elangreza/e-commerce/payment/internal/sqlitedb"
	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ServicePort string `koanf:"SERVICE_PORT"`
}

func main() {
	var cfg Config
	err := config.LoadConfig(&cfg)
	errChecker(err)

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("payment.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)

	paymentRepo := sqlitedb.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	srv := server.New(paymentService)

	addr := fmt.Sprintf(":%s", cfg.ServicePort)

	go func() {
		if err := srv.Start(addr); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	fmt.Printf("MOCKED-PAYMENT-service running at %s\n", addr)

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
