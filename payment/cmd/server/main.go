package main

import (
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"

	"log"

	"github.com/elangreza/e-commerce/payment/internal/server"
	"github.com/elangreza/e-commerce/payment/internal/service"
	"github.com/elangreza/e-commerce/payment/internal/sqlitedb"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("payment.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)
	defer db.Close()

	paymentRepo := sqlitedb.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	srv := server.New(paymentService)

	address := fmt.Sprintf(":%v", 50053)

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
