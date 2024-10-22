package handler

import (
	"net/http"

	"hot-coffee/internal/service"
)

type AggregationHandler struct {
	orderService *service.OrderService
}

func NewAggregationHandler(orderService *service.OrderService) *AggregationHandler {
	return &AggregationHandler{orderService: orderService}
}

// Return all saled items as key and quantity as value in JSON
func (h *AggregationHandler) TotalSalesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// TODO
		return
	}
}

// Returns Each item as key and quatity as value
func (h *AggregationHandler) PopularItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// TODO
		return
	}
}
