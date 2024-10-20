package OrdersHandlers

import (
	"encoding/json"
	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/services"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func PostOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	if !services.MenuCheck(w, order) {
		ErrorHandler.Error(w, "Your order does not exist in menu", http.StatusBadRequest)
		return
	}
	if !services.IngredientsCheck(w, order) {
		ErrorHandler.Error(w, "Not enough ingedients or needed ingerdients do not exist", http.StatusBadRequest)
		return
	}

	var orders []models.Order
	content, err := ioutil.ReadFile(config.BaseDir + "/orders.json") // ioutil.ReadFile() читает файл и возвращает содержимое в массие из байтов
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(content, &orders) // json.Unmarshal([]byte, type any) короче из инфы в байтах он конвертирует все в структуру в стиле json

	Location, err := time.LoadLocation("Asia/Aqtau")
	timenow := time.Now().In(Location).Format(time.RFC3339)
	order.CreatedAt = timenow
	order.ID = strconv.Itoa(GetID(w))
	order.Status = "open"

	orders = append(orders, order)

	// json.MarshalIndent() принимает структуру и конвертирует все в инфу в стиле json
	jsonData, err := json.MarshalIndent(orders, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could not convert orders to json data", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(config.BaseDir+"/orders.json", jsonData, 0644) // os.WriteFile(filename, content, perm) в файл записывает данные
	if err != nil {
		ErrorHandler.Error(w, "Could not write orders to json database", http.StatusInternalServerError)
		return
	}
}

func GetID(w http.ResponseWriter) int {
	ConfigContent, err := os.ReadFile(config.BaseDir + "/config.json")
	if err != nil {
		ErrorHandler.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	var ID models.OrderID
	err = json.Unmarshal(ConfigContent, &ID)
	if err != nil {
		ErrorHandler.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	i := ID.ID
	ID.ID++
	NewContent, err := json.MarshalIndent(ID, "", "    ")
	os.WriteFile(config.BaseDir+"/config.json", NewContent, os.ModePerm)
	return i
}
