package handler

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/constanta"
	"github.com/elangreza/e-commerce/pkg/money"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	svc  gen.PaymentServiceServer
	tmpl *template.Template
}

func NewHandler(
	tmpl *template.Template,
	publicRoute chi.Router,
	svc gen.PaymentServiceServer,
) {

	h := handler{
		tmpl: tmpl,
		svc:  svc,
	}

	publicRoute.Get("/", h.indexGet)
	publicRoute.Post("/", h.indexPost)

	publicRoute.Get("/transactions/{transactionID}", h.detailGet)
	publicRoute.Post("/transactions/{transactionID}", h.detailPost)

	publicRoute.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.tmpl.ExecuteTemplate(w, "404.html", nil)
	})
}

type PageData struct {
	// Shared
	Error string

	// Detail
	Transaction *Transaction
}

type Transaction struct {
	ID          string
	TotalAmount string
	Status      constanta.PaymentStatus
	ExpiredAt   string
	CreatedAt   string
}

func (h *handler) indexGet(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (h *handler) indexPost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("transaction_id"))
	if id == "" {
		h.renderError(w, "index.html", "Invalid transaction ID")
		return
	}
	http.Redirect(w, r, "/transactions/"+id, http.StatusSeeOther)
}

// --- Detail Handlers ---

func (h *handler) detailGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "transactionID")
	if id == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	payment, err := h.svc.GetPayment(context.Background(), &gen.GetPaymentRequest{TransactionId: id})
	if err != nil {
		h.handleDetailError(w, id, err)
		return
	}

	total, err := money.ToMajorString(payment.TotalAmount)
	if err != nil {
		h.renderError(w, "detail.html", &PageData{
			Transaction: &Transaction{ID: id, Status: "UNKNOWN"},
			Error:       "Invalid transaction amount",
		})
		return
	}

	data := PageData{
		Transaction: &Transaction{
			ID:          id,
			TotalAmount: total,
			Status:      constanta.PaymentStatus(payment.Status),
			CreatedAt:   payment.CreatedAt,
			ExpiredAt:   payment.ExpiredAt,
		},
	}

	h.tmpl.ExecuteTemplate(w, "detail.html", data)
}

func (h *handler) detailPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "transactionID")
	if id == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	// Re-fetch transaction to check current status
	payment, err := h.svc.GetPayment(context.Background(), &gen.GetPaymentRequest{TransactionId: id})
	if err != nil {
		h.handleDetailError(w, id, err)
		return
	}

	// Only allow payment if status is WAITING
	if constanta.PaymentStatus(payment.Status) != constanta.WAITING {
		http.Redirect(w, r, "/transactions/"+id, http.StatusSeeOther)
		return
	}

	paymentStr := strings.TrimSpace(r.FormValue("payment_amount"))
	amount, err := strconv.Atoi(paymentStr)
	if err != nil {
		// Re-render with error
		total, _ := money.ToMajorString(payment.TotalAmount) // safe since already validated in GET
		data := PageData{
			Transaction: &Transaction{
				ID:          id,
				TotalAmount: total,
				Status:      constanta.WAITING,
			},
			Error: "Invalid payment amount",
		}
		h.tmpl.ExecuteTemplate(w, "detail.html", data)
		return
	}

	_, err = h.svc.UpdatePayment(context.Background(), &gen.UpdatePaymentRequest{
		TransactionId: id,
		TotalAmount:   &gen.Money{Units: int64(amount), CurrencyCode: "IDR"},
	})

	if err != nil {
		// Re-render with error
		total, _ := money.ToMajorString(payment.TotalAmount)
		data := PageData{
			Transaction: &Transaction{
				ID:          id,
				TotalAmount: total,
				Status:      constanta.WAITING,
			},
			Error: "Failed to process payment",
		}
		h.tmpl.ExecuteTemplate(w, "detail.html", data)
		return
	}

	// Redirect to same detail page (now showing updated status)
	http.Redirect(w, r, "/transactions/"+id, http.StatusSeeOther)
}

// --- Helpers ---

func (h *handler) handleDetailError(w http.ResponseWriter, id string, err error) {
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		h.renderError(w, "detail.html", &PageData{
			Transaction: &Transaction{ID: id, Status: "UNKNOWN"},
			Error:       "Transaction not found",
		})
		return
	}

	// For any other error
	h.renderError(w, "detail.html", &PageData{
		Transaction: &Transaction{ID: id, Status: "UNKNOWN"},
		Error:       "Unable to load transaction",
	})
}

func (h *handler) renderError(w http.ResponseWriter, tmplName string, data interface{}) {
	h.tmpl.ExecuteTemplate(w, tmplName, data)
}
