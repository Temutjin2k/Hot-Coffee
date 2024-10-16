package handler

import (
	"io/ioutil"
	"net/http"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("data/orders.json")
	if err != nil {
		// will be handler
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
