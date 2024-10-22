package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
	logger           *slog.Logger
}

func NewInventoryHandler(inventoryService *service.InventoryService, logger *slog.Logger) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService, logger: logger}
}

func (h *InventoryHandler) PostInventory(w http.ResponseWriter, r *http.Request) {
	var newItem models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	err = h.inventoryService.AddInventoryItem(newItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not add new inventory item", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) PutInventoryItem(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
}
