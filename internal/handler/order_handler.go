package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/orderHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[1:], "/")
	switch len(parts) {
	case 1:
		switch r.Method {
		case http.MethodPost:
			h.PostOrder(w, r)
		case http.MethodGet:
			h.GetOrders(w, r)
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 2:
		switch r.Method {
		case http.MethodPut:
			h.PutOrder(w, r)
		case http.MethodGet:
			h.GetOrder(w, r)
		case http.MethodDelete:
			h.DeleteOrder(w, r)
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 3:
		if r.Method == http.MethodPost {
			if parts[2] == "close" {
				orderHandler.Closeorder(w, parts[1])
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
	json.NewEncoder(w).Encode(newOrder)
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
