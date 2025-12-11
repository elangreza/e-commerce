package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/elangreza/e-commerce/payment/internal/handler"
	"github.com/elangreza/e-commerce/payment/internal/server"
	"github.com/elangreza/e-commerce/payment/internal/service"
	"github.com/elangreza/e-commerce/payment/internal/sqlitedb"
	"github.com/elangreza/e-commerce/pkg/config"
	"github.com/elangreza/e-commerce/pkg/dbsql"
	"github.com/elangreza/e-commerce/pkg/gracefulshutdown"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	ServicePort        string        `koanf:"SERVICE_PORT"`
	MaxTimeToBeExpired time.Duration `koanf:"MAX_TIME_TO_BE_EXPIRED"`
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
	paymentService := service.NewPaymentService(paymentRepo, cfg.MaxTimeToBeExpired)
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
	)
	<-gs
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// var tmpl *template.Template

// type PageData struct {
// 	// Shared
// 	Error string

// 	// Index
// 	TransactionID string

// 	// Detail
// 	Transaction   *Transaction
// 	PaymentAmount string

// 	// Status
// 	FinalStatus string
// }

// type Transaction struct {
// 	ID          string
// 	TotalAmount float64
// 	Status      constanta.PaymentStatus
// 	IsExpired   bool
// 	ExpiredAt   string
// }

// // 1. INDEX PAGE: Enter transaction ID
// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		idStr := strings.TrimSpace(r.FormValue("transaction_id"))
// 		if idStr == "" {
// 			tmpl.ExecuteTemplate(w, "index.html", PageData{Error: "Invalid transaction ID"})
// 			return
// 		}

// 		// Redirect to detail page
// 		http.Redirect(w, r, "/transaction/"+idStr, http.StatusSeeOther)
// 		return
// 	}

// 	tmpl.ExecuteTemplate(w, "index.html", PageData{})
// }

// // 2. TRANSACTION DETAIL PAGE
// func detailHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extract ID from URL path
// 	idStr := strings.TrimPrefix(r.URL.Path, "/transaction/")
// 	if idStr == "" {
// 		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch transaction
// 	var t = Transaction{
// 		ID:          idStr,
// 		TotalAmount: 10000,
// 		Status:      constanta.WAITING,
// 		ExpiredAt:   time.Now().Add(3 * time.Minute).Format(time.DateTime),
// 	}

// 	if r.Method == "POST" {
// 		if t.Status != "pending" {
// 			// Already finalized – redirect to status
// 			http.Redirect(w, r, "/status/"+idStr, http.StatusSeeOther)
// 			return
// 		}

// 		paymentStr := strings.TrimSpace(r.FormValue("payment_amount"))
// 		_, err := strconv.ParseFloat(paymentStr, 64)
// 		if err != nil {
// 			data := PageData{
// 				Transaction:   &t,
// 				PaymentAmount: paymentStr,
// 				Error:         "Invalid payment amount",
// 			}
// 			tmpl.ExecuteTemplate(w, "detail.html", data)
// 			return
// 		}

// 		return
// 	}

// 	// Render detail page
// 	data := PageData{
// 		Transaction: &t,
// 	}
// 	tmpl.ExecuteTemplate(w, "detail.html", data)
// }
