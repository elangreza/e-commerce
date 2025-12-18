package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/handler"
	"github.com/elangreza/e-commerce/payment/internal/server"
	"github.com/elangreza/e-commerce/payment/internal/service"
	"github.com/elangreza/e-commerce/payment/internal/sqlitedb"
	"github.com/elangreza/e-commerce/payment/task"
	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ServicePort        string        `koanf:"SERVICE_PORT"`
	DBPath             string        `koanf:"DB_PATH"`
	MaxTimeToBeExpired time.Duration `koanf:"MAX_TIME_TO_BE_EXPIRED"`
	OrderServiceAddr   string        `koanf:"ORDER_SERVICE_ADDR"`
}

func main() {
	var cfg Config
	err := config.LoadConfig(&cfg)
	errChecker(err)

	// Default to data-local for local development
	dbPath := cfg.DBPath
	if dbPath == "" {
		dbPath = "data-local/payment.db"
	}

	// implement this later
	// github.com/samber/slog-zap

	db, err := dbsql.NewDbSql(
		dbsql.WithSqliteDB(dbPath),
		dbsql.WithSqliteDBWalMode(),
		dbsql.WithAutoMigrate("file://./migrations"),
	)
	errChecker(err)

	// order client
	grpcClientOrder, err := grpc.NewClient(cfg.OrderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	errChecker(err)

	paymentRepo := sqlitedb.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cfg.MaxTimeToBeExpired, gen.NewOrderServiceClient(grpcClientOrder))
	srv := server.New(paymentService)

	addr := fmt.Sprintf(":%s", cfg.ServicePort)

	go func() {
		if err := srv.Start(addr); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	h := chi.NewRouter()
	h.Use(middleware.Recoverer)
	h.Use(middleware.Logger)
	h.Use(middleware.Timeout(60 * time.Second))
	h.Use(middleware.RequestID)
	h.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	handler.NewHandler(tmpl, h, paymentService)
	srvHttp := &http.Server{
		Addr:           ":8081",
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srvHttp.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	taskPayment := task.NewTaskPayment(paymentService, cfg.MaxTimeToBeExpired)

	fmt.Printf("MOCKED-PAYMENT-service running at %s\n", addr)
	fmt.Println("UI-MOCKED-PAYMENT-service running on :8081")

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
		gracefulshutdown.Operation{
			Name: "server",
			ShutdownFunc: func(ctx context.Context) error {
				return srvHttp.Shutdown(ctx)
			},
		},
		gracefulshutdown.Operation{
			Name: "task payment",
			ShutdownFunc: func(ctx context.Context) error {
				taskPayment.Close()
				return nil
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
