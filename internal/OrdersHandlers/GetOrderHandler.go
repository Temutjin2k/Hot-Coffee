package OrdersHandlers

import (
	"encoding/json"
	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
)

func GetOrders(w http.ResponseWriter) {
	content, err := ioutil.ReadFile(config.BaseDir + "/orders.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func GetOrder(w http.ResponseWriter, OrderID string) {
	OrderContents, err := ioutil.ReadFile(config.BaseDir + "/orders.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	flag := true
	var Orders []models.Order
	var NeededOrder models.Order
	json.Unmarshal(OrderContents, &Orders)
	for _, order := range Orders {
		if order.ID == OrderID {
			NeededOrder.CreatedAt = order.CreatedAt
			NeededOrder.CustomerName = order.CustomerName
			NeededOrder.ID = OrderID
			NeededOrder.Items = order.Items
			NeededOrder.Status = order.Status
			flag = false
		}
	}
	jsondata, err := json.MarshalIndent(NeededOrder, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could not upload order", http.StatusInternalServerError)
		return
	}
	if !flag {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsondata)
	} else {
		ErrorHandler.Error(w, "Your requested order is not found", 404)
	}
}
