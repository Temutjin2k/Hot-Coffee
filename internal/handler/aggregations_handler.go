package handler

import "net/http"

type AggregationHandler struct{}

func NewAggregationHandler() *AggregationHandler {
	return &AggregationHandler{}
}

func (h *AggregationHandler) TotalSalesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// TODO
	}
}

func (h *AggregationHandler) PopularItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// TODO
	}
}
