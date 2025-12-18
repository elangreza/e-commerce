package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elangreza/e-commerce/pkg/config"
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

type Config struct {
	ServicePort          string `koanf:"SERVICE_PORT"`
	TokenSecret          string `koanf:"TOKEN_SECRET"`
	DBPath               string `koanf:"DB_PATH"`
	OrderServiceAddr     string `koanf:"ORDER_SERVICE_ADDR"`
	ProductServiceAddr   string `koanf:"PRODUCT_SERVICE_ADDR"`
	WarehouseServiceAddr string `koanf:"WAREHOUSE_SERVICE_ADDR"`
	ShopServiceAddr      string `koanf:"SHOP_SERVICE_ADDR"`
}

func main() {
	var cfg Config
	err := config.LoadConfig(&cfg)
	errChecker(err)

	// Default to data-local for local development
	dbPath := cfg.DBPath
	if dbPath == "" {
		dbPath = "data-local/auth.db"
	}

	handler := chi.NewRouter()
	handler.Use(middleware.Recoverer)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Timeout(60 * time.Second))
	handler.Use(middleware.RequestID)
	handler.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB(dbPath),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)

	// repositories
	userRepo := sqlitedb.NewUserRepo(db)
	tokenRepo := sqlitedb.NewTokenRepo(db)

	// order
	grpcClientOrder, err := grpc.NewClient(cfg.OrderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// product
	grpcClientProduct, err := grpc.NewClient(cfg.ProductServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// warehouse
	grpcClientWarehouse, err := grpc.NewClient(cfg.WarehouseServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// shop
	grpcClientShop, err := grpc.NewClient(cfg.ShopServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	// services
	authService := service.NewAuthService(userRepo, tokenRepo, cfg.TokenSecret)
	productService := service.NewProductService(gen.NewProductServiceClient(grpcClientProduct), gen.NewShopServiceClient(grpcClientShop))
	orderService := service.NewOrderService(gen.NewOrderServiceClient(grpcClientOrder))
	warehouseService := service.NewWarehouseService(gen.NewWarehouseServiceClient(grpcClientWarehouse))

	rest.NewAuthHandler(handler, authService)
	rest.NewProductHandler(handler, productService)
	rest.NewOrderHandler(handler, authService, orderService)
	rest.NewWarehouseHandler(handler, authService, warehouseService)

	addr := fmt.Sprintf(":%s", cfg.ServicePort)

	srv := &http.Server{
		Addr:           addr,
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

	fmt.Printf("API-service running at %s\n", addr)

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
