package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

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

func (h *OrderHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[1:], "/")
	switch len(parts) {
	case 1: // Endpoint: /orders
		switch r.Method {
		case http.MethodPost:
			h.PostOrder(w, r) // POST /orders: Create a new order.
		case http.MethodGet:
			h.GetOrders(w, r) // GET /orders: Retrieve all orders.
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 2: // Endpoint: /orders/{id}
		switch r.Method {
		case http.MethodPut:
			h.PutOrder(w, r) // PUT /orders/{id}: Update an existing order.
		case http.MethodGet:
			h.GetOrder(w, r) // GET /orders/{id}: Retrieve a specific order by ID.
		case http.MethodDelete:
			h.DeleteOrder(w, r) // DELETE /orders/{id}: Delete an order.
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 3: // Endpoint: /orders/{id}/close
		if r.Method == http.MethodPost {
			if parts[2] == "close" {
				h.CloseOrder(w, r) // POST /orders/{id}/close: Close an order.
			} else {
				ErrorHandler.Error(w, "Adress is not allowed", http.StatusForbidden)
			}
		} else {
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	}
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
