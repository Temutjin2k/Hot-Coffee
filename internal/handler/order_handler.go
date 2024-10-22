package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type OrderHandler struct {
	orderService *service.OrderService
	logger       *slog.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, logger: logger}
}

// PostOrder creates new Order
func (h *OrderHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	err = h.orderService.AddOrder(newOrder)
	if err != nil {
		ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrder) // TODO error handling
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
}

func (h *OrderHandler) PutOrder(w http.ResponseWriter, r *http.Request) {
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
}
