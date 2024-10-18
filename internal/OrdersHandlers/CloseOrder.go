package OrdersHandlers

import (
	"encoding/json"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/services"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
)

func Closeorder(w http.ResponseWriter, OrderID string) {
	OrderContents, err := ioutil.ReadFile("data/orders.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	var Orders []models.Order
	json.Unmarshal(OrderContents, &Orders)

	var ClosingOrder models.Order
	for _, order := range Orders {
		if order.ID == OrderID {
			ClosingOrder.CreatedAt = order.CreatedAt
			ClosingOrder.CustomerName = order.CustomerName
			ClosingOrder.ID = OrderID
			ClosingOrder.Items = order.Items
			ClosingOrder.Status = "Closed"
		}
	}
	for i, order := range Orders {
		if order.ID == OrderID {
			Orders[i].CreatedAt = ClosingOrder.CreatedAt
			Orders[i].CustomerName = ClosingOrder.CustomerName
			Orders[i].ID = OrderID
			Orders[i].Items = ClosingOrder.Items
			Orders[i].Status = "Closed"
		}
	}

	services.SubtractIngridients(w, ClosingOrder)
	jsondata, err := json.MarshalIndent(Orders, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could not upload order", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile("data/orders.json", jsondata, os.ModePerm)
	if err != nil {
		ErrorHandler.Error(w, "Could not rewrite orders", http.StatusInternalServerError)
		return
	}
}
