package handler

import (
	"fmt"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/MenuHandlers"
	"hot-coffee/internal/OrdersHandlers"
	"net/http"
	"strings"
)

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	Parts := strings.SplitN(r.URL.Path[1:], "/", 2)
	fmt.Println(Parts)
	fmt.Println(len(Parts))
	switch len(Parts) {
	case 1:
		switch r.Method {
		case http.MethodPost:
			MenuHandlers.MenuPost(w, r)
		case http.MethodGet:
			MenuHandlers.GetMenuItems(w)
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	case 2:
		switch r.Method {
		case http.MethodPut:
			MenuHandlers.MenuPut(w, r, Parts[1])
		case http.MethodGet:
			MenuHandlers.GetMenuItem(w, Parts[1])
		case http.MethodDelete:
			OrdersHandlers.Deleteorder(w, Parts[1])
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	}
}
