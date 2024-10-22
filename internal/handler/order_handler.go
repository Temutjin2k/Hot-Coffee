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
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	for _, OrderItem := range NewOrder.Items {
		if !h.menuService.MenuCheckByID(OrderItem.ProductID) {
			ErrorHandler.Error(w, "Updated order item does not exist in menu", http.StatusBadRequest)
			return
		}
		if !h.menuService.IngredientsCheckByID(OrderItem.ProductID, OrderItem.Quantity) {
			ErrorHandler.Error(w, "Not enough ingridients for your upddated order", http.StatusBadRequest)
			return
		}
	}

	err = h.orderService.AddOrder(NewOrder)
	if err != nil {
		ErrorHandler.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	RequestedOrder, err := h.orderService.GetOrder(r.PathValue("id"))
	if err != nil {
		ErrorHandler.Error(w, "Something happened when getting order", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.MarshalIndent(RequestedOrder, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Can not convert order data to json", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	Orders, err := h.orderService.GetAllOrders()
	if err != nil {
		ErrorHandler.Error(w, "Can not read order data from server", http.StatusInternalServerError)
	}
	jsonData, err := json.MarshalIndent(Orders, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Can not convert order data to json", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func (h *OrderHandler) PutOrder(w http.ResponseWriter, r *http.Request) {
	Requestcontent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}
	var RequestedOrder models.Order
	err = json.Unmarshal(Requestcontent, &RequestedOrder)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	for _, OrderItem := range RequestedOrder.Items {
		if !h.menuService.MenuCheckByID(OrderItem.ProductID) {
			ErrorHandler.Error(w, "Updated order item does not exist in menu", http.StatusBadRequest)
			return
		}
		if !h.menuService.IngredientsCheckByID(OrderItem.ProductID, OrderItem.Quantity) {
			ErrorHandler.Error(w, "Not enough ingridients for your upddated order", http.StatusBadRequest)
			return
		}
	}
	h.orderService.UpdateOrder(RequestedOrder, r.PathValue("id"))
	w.WriteHeader(200)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	err := h.orderService.SaveAllOrders(r.PathValue("id"))
	if err != nil {
		ErrorHandler.Error(w, "Error updating orders database", http.StatusBadRequest)
	}
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	Order, err := h.orderService.GetOrder(r.PathValue("id"))
	if err != nil {
		ErrorHandler.Error(w, "Something happened when getting order", http.StatusInternalServerError)
		return
	}
	for _, item := range Order.Items {
		err := h.menuService.SubtractIngredientsByID(item.ProductID, item.Quantity)
		if err != nil {
			ErrorHandler.Error(w, "Not enough ingridients to close the order", http.StatusInternalServerError)
		}
	}
	err = h.orderService.CloseOrder(r.PathValue("id"))
	if err != nil {
		ErrorHandler.Error(w, "Something happened when closing order", http.StatusInternalServerError)
		return
	}
}
