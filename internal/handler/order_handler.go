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

type OrderHandler struct {
	orderService *service.OrderService
	menuService  *service.MenuService
	logger       *slog.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService, menuService *service.MenuService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, menuService: menuService, logger: logger}
}

// PostOrder creates new Order
func (h *OrderHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var NewOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&NewOrder)
	if err != nil {
		h.logger.Error("Could not decode request json data", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	for _, OrderItem := range NewOrder.Items {
		err = h.menuService.MenuCheckByID(OrderItem.ProductID)
		if err != nil {
			h.logger.Error("Requested order item does not exist in menu", "error", err, "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, "Requested order item does not exist in menu", http.StatusBadRequest)
			return
		}
		if err = h.menuService.IngredientsCheckByID(OrderItem.ProductID, OrderItem.Quantity); err != nil {
			h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = h.orderService.AddOrder(NewOrder)
	if err != nil {
		h.logger.Error("Something wrong when adding new order", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Something wrong when adding new order", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	RequestedOrder, err := h.orderService.GetOrder(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Something wrong when getting order by ID", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Something wrong when getting order by ID", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.MarshalIndent(RequestedOrder, "", "    ")
	if err != nil {
		h.logger.Error("Can not convert order data to json", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Can not convert order data to json", http.StatusInternalServerError)
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	Orders, err := h.orderService.GetAllOrders()
	if err != nil {
		h.logger.Error("Can not read order data from server", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Can not read order data from server", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(Orders, "", "    ")
	if err != nil {
		h.logger.Error("Can not convert order data to json", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Can not convert order data to json", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
}

func (h *OrderHandler) PutOrder(w http.ResponseWriter, r *http.Request) {
	Requestcontent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Could not read request body", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}
	var RequestedOrder models.Order
	err = json.Unmarshal(Requestcontent, &RequestedOrder)
	if err != nil {
		h.logger.Error("Could not read request body", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	for _, OrderItem := range RequestedOrder.Items {
		err = h.menuService.MenuCheckByID(OrderItem.ProductID)
		if err != nil {
			h.logger.Error("Updated order item does not exist in menu", "error", err, "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, "Updated order item does not exist in menu", http.StatusBadRequest)
			return
		}
		if err = h.menuService.IngredientsCheckByID(OrderItem.ProductID, OrderItem.Quantity); err != nil {
			h.logger.Error(err.Error(), "error", err, "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	err = h.orderService.UpdateOrder(RequestedOrder, r.PathValue("id"))
	if err != nil {
		h.logger.Error("Could not update order", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Could not update order", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.WriteHeader(200)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	err := h.orderService.DeleteOrderByID(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Error updating orders database", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Error updating orders database", http.StatusInternalServerError)
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.WriteHeader(204)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	Order, err := h.orderService.GetOrder(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Something happened when getting order by ID", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Something happened when getting order by ID", http.StatusInternalServerError)
		return
	}
	for _, item := range Order.Items {
		err := h.menuService.SubtractIngredientsByID(item.ProductID, item.Quantity)
		if err != nil {
			h.logger.Error("Not enough ingridients to close the order", "error", err, "method", r.Method, "url", r.URL)
			ErrorHandler.Error(w, "Not enough ingridients to close the order", http.StatusBadRequest)
		}
	}
	err = h.orderService.CloseOrder(r.PathValue("id"))
	if err != nil {
		h.logger.Error("Something happened when closing order", "error", err, "method", r.Method, "url", r.URL)
		ErrorHandler.Error(w, "Something happened when closing order", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Request handled successfully.", "method", r.Method, "url", r.URL)
	w.WriteHeader(200)
}
