package handler

import (
	"encoding/json"
	"hot-coffee/internal/businesslogic"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func PostOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		// will write error handler
	}
	businesslogic.IngredientsCheck(order)

	var orders []models.Order
	content, err := ioutil.ReadFile("data/orders.json") // ioutil.ReadFile() читает файл и возвращает содержимое в массие из байтов
	if err != nil {
		// will write error handler
	}
	json.Unmarshal(content, &orders) // json.Unmarshal([]byte, type any) короче из инфы в байтах он конвертирует все в структуру в стиле json

	Location, err := time.LoadLocation("Asia/Aqtau")
	timenow := time.Now().In(Location).Format(time.RFC3339)
	order.ID = "1"
	order.Status = "active"
	order.CreatedAt = timenow

	orders = append(orders, order)

	// json.MarshalIndent() принимает структуру и конвертирует все в инфу в стиле json
	jsonData, err := json.MarshalIndent(orders, "", "    ")
	if err != nil {
		// will write error handler
	}

	err = os.WriteFile("data/orders.json", jsonData, 0644) // os.WriteFile(filename, content, perm) в файл записывает данные
	if err != nil {
		// will write error handler
	}
}
