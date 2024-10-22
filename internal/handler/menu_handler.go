package handler

import (
	"encoding/json"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
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
	MenuItems, err := h.menuService.GetMenuItems()
	if err != nil {
		ErrorHandler.Error(w, "Could not read menu database", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(MenuItems, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	MenuItem, err := h.menuService.GetMenuItem(r.PathValue("id"))
	if err != nil {
		ErrorHandler.Error(w, "Could not read menu database", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(MenuItem, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *MenuHandler) PutMenuItem(w http.ResponseWriter, r *http.Request) {
	err := h.menuService.UpdateMenuItem(r)
	if err != nil {
		h.logger.Error("Could not update menu database", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not update menu database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	err := h.menuService.DeleteMenuItem(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Could not delete menu item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not delete menu item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(204)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}
