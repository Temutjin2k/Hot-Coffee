package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type MenuHandler struct {
	menuService *service.MenuService
	logger      *slog.Logger
}

func NewMenuHandler(menuService *service.MenuService, logger *slog.Logger) *MenuHandler {
	return &MenuHandler{menuService: menuService, logger: logger}
}

func (h *MenuHandler) PostMenu(w http.ResponseWriter, r *http.Request) {
	var newItem models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	// Use the service to check if the item already exists
	if h.menuService.MenuCheck(newItem) {
		ErrorHandler.Error(w, "The requested menu item already exists in current menu", http.StatusBadRequest)
		return
	}

	// Add the new menu item using the service
	if err := h.menuService.AddMenuItem(newItem); err != nil {
		ErrorHandler.Error(w, "Could not add menu item", http.StatusInternalServerError)
		return
	}
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
}

func (h *MenuHandler) PutMenuItem(w http.ResponseWriter, r *http.Request) {
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
}
