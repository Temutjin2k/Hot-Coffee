package handler

import (
	"net/http"
	"strings"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) InventoryHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path[1:], "/")

	switch len(path) {
	case 1: // Endpoint: /inventory
		switch r.Method {
		case http.MethodPost: // POST /inventory: Add a new inventory item.
			h.PostInventory(w, r)
		case http.MethodGet: // GET /inventory: Retrieve all inventory items.
			h.GetInventory(w, r)
		}
	case 2: // Endpoint: /inventory/{id}
		// Maybe Validation of ID
		switch r.Method {
		case http.MethodGet: // GET /inventory/{id}: Retrieve a specific inventory item.
			h.GetInventoryItem(w, r)
		case http.MethodPut: // PUT /inventory/{id}: Update an inventory item.
			h.PutInventoryItem(w, r)
		case http.MethodDelete: // DELETE /inventory/{id}: Delete an inventory item.
			h.DeleteInventoryItem(w, r)
		}
	default:
		ErrorHandler.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (h *InventoryHandler) PostInventory(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) PutInventoryItem(w http.ResponseWriter, r *http.Request) {
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
}
