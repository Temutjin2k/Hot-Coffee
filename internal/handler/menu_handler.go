package handler

import (
	"encoding/json"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"io/ioutil"
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
		h.logger.Error("Could not decode request json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}
	err = h.menuService.CheckNewMenu(newItem)
	if err != nil {
		h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use the service to check if the item already exists
	if err = h.menuService.MenuCheckByID(newItem.ID); err != nil {
		h.logger.Error("The requested menu item already exists in current menu", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "The requested menu item already exists in current menu", http.StatusBadRequest)
		return
	}
	if err = h.menuService.IngredientsCheckByID(newItem.ID, 1); err != nil {
		h.logger.Error(err.Error(), "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the new menu item using the service
	if err := h.menuService.AddMenuItem(newItem); err != nil {
		h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not add menu item", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	MenuItems, err := h.menuService.GetMenuItems()
	if err != nil {
		h.logger.Error("Could not read menu database", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not read menu database", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(MenuItems, "", "    ")
	if err != nil {
		h.logger.Error("Could not read menu database", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not read menu items", http.StatusInternalServerError)
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	MenuItem, err := h.menuService.GetMenuItem(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Could not get menu item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not read menu database", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(MenuItem, "", "    ")
	if err != nil {
		h.logger.Error("Could not convert Menu Items to jsondata", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not send menu item", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)

	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *MenuHandler) PutMenuItem(w http.ResponseWriter, r *http.Request) {
	err := h.menuService.MenuCheckByID(r.PathValue("id"))
	if err != nil {
		h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	RequestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error reading body of request", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Something wrong with your request", http.StatusInternalServerError)
		return
	}

	var RequestedMenuItem models.MenuItem
	err = json.Unmarshal(RequestContent, &RequestedMenuItem)
	if err != nil {
		h.logger.Error("Error converting menu item to json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Something wrong with your request", http.StatusInternalServerError)
		return
	}
	err = h.menuService.CheckNewMenu(RequestedMenuItem)
	if err != nil {
		h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = h.menuService.IngredientsCheckByID(RequestedMenuItem.ID, 1); err != nil {
		h.logger.Error(err.Error(), "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.menuService.UpdateMenuItem(r)
	if err != nil {
		h.logger.Error("Could not update menu database", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not update menu database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	err := h.menuService.MenuCheckByID(r.PathValue("id"))
	if err != nil {
		h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.menuService.DeleteMenuItem(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Could not delete menu item", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not delete menu item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(204)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}
