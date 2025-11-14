package main

import (
	"fmt"
	"github/elangreza/e-commerce/pkg/dbsql"
	"github/elangreza/e-commerce/pkg/interceptor"

	"log"
	"net"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/grpcserver"
	"github.com/elangreza/e-commerce/payment/internal/service"
	"github.com/elangreza/e-commerce/payment/internal/sqlitedb"
	"google.golang.org/grpc"
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
	paymentServer := grpcserver.NewPaymentServer(paymentService)

	address := fmt.Sprintf(":%v", 50053)
	listener, err := net.Listen("tcp", address)
	errChecker(err)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.UserIDParser(),
		),
	)
	gen.RegisterPaymentServiceServer(grpcServer, paymentServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
