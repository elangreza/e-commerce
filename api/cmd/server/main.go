package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"

	"github.com/elangreza/e-commerce/api/internal/rest"
	"github.com/elangreza/e-commerce/api/internal/service"
	"github.com/elangreza/e-commerce/api/internal/sqlitedb"

	"github.com/elangreza/e-commerce/gen"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	handler := chi.NewRouter()
	handler.Use(middleware.Recoverer)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Timeout(60 * time.Second))
	handler.Use(middleware.RequestID)
	handler.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB("auth.db"),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)

	// repositories
	userRepo := sqlitedb.NewUserRepo(db)
	tokenRepo := sqlitedb.NewTokenRepo(db)

	// service
	authService := service.NewAuthService(userRepo, tokenRepo)

	// order
	grpcClientOrder, err := grpc.NewClient("order:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// product
	grpcClientProduct, err := grpc.NewClient("product:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// warehouse
	grpcClientWarehouse, err := grpc.NewClient("warehouse:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// shop
	grpcClientShop, err := grpc.NewClient("shop:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	productService := service.NewProductService(gen.NewProductServiceClient(grpcClientProduct), gen.NewShopServiceClient(grpcClientShop))
	orderService := service.NewOrderService(gen.NewOrderServiceClient(grpcClientOrder))
	warehouseService := service.NewWarehouseService(gen.NewWarehouseServiceClient(grpcClientWarehouse))

	rest.NewAuthHandler(handler, authService)
	rest.NewProductHandler(handler, productService)
	rest.NewOrderHandler(handler, authService, orderService)
	rest.NewWarehouseHandler(handler, authService, warehouseService)

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	gs := gracefulshutdown.New(context.Background(), 5*time.Second,
		gracefulshutdown.Operation{
			Name: "server",
			ShutdownFunc: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			}},
		gracefulshutdown.Operation{
			Name: "sqlite",
			ShutdownFunc: func(ctx context.Context) error {
				return db.Close()
			}},
	)
	<-gs
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
