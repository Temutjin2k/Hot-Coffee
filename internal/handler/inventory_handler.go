package handler

import (
	"net/http"

	"hot-coffee/internal/service"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) InventoryHandler(w http.ResponseWriter, r *http.Request) {
}
