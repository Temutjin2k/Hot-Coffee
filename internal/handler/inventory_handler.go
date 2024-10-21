package handler

import "net/http"

type InventoryHandler struct{}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

func (h *InventoryHandler) InventoryHandler(w http.ResponseWriter, r *http.Request) {
}
