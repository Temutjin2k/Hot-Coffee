package Mainhandlers

import (
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/OrdersHandlers"
	"net/http"
	"strings"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	Parts := strings.SplitN(r.URL.Path[1:], "/", 3)
	switch len(Parts) {
	case 1:
		switch r.Method {
		case http.MethodPost:
			OrdersHandlers.PostOrder(w, r)
		case http.MethodGet:
			OrdersHandlers.GetOrders(w)
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 2:
		switch r.Method {
		case http.MethodPut:
			OrdersHandlers.Putorder(w, r, Parts[1])
		case http.MethodGet:
			OrdersHandlers.GetOrder(w, Parts[1])
		case http.MethodDelete:
			OrdersHandlers.Deleteorder(w, Parts[1])
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 3:
		if r.Method == http.MethodPost {
			if Parts[2] == "close" {
				OrdersHandlers.Closeorder(w, Parts[1])
			} else {
				ErrorHandler.Error(w, "Adress is not allowed", http.StatusForbidden)
			}
		} else {
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	}
}
