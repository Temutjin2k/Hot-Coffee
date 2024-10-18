package OrdersHandlers

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/services"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
)

func Putorder(w http.ResponseWriter, r *http.Request, OrderID string) {
	Requestcontent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var RequestrOrder models.Order
	err = json.Unmarshal(Requestcontent, &RequestrOrder)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	if !services.IngredientsCheck(w, RequestrOrder) {
		ErrorHandler.Error(w, "Not enough ingridients of your new order or theu do not exist ", http.StatusBadRequest)
		return
	}

	OrderContents, err := ioutil.ReadFile("data/orders.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	var Orders []models.Order
	json.Unmarshal(OrderContents, &Orders)

	for i, order := range Orders {
		if order.ID == OrderID {
			fmt.Println(OrderID)
			Orders[i].CreatedAt = RequestrOrder.CreatedAt
			Orders[i].CustomerName = RequestrOrder.CustomerName
			Orders[i].ID = OrderID
			Orders[i].Items = RequestrOrder.Items
			Orders[i].Status = RequestrOrder.Status
		}
	}

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
