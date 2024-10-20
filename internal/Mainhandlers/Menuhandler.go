package Mainhandlers

import (
	"fmt"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/MenuHandlers"
	"net/http"
	"strings"
)

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	Parts := strings.Split(r.URL.Path[1:], "/")
	switch len(Parts) {
	case 1:
		fmt.Println(r.Method)
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
			MenuHandlers.MenuDelete(w, Parts[1])
		default:
			ErrorHandler.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	default:
		ErrorHandler.Error(w, "Something wrong with your request", http.StatusBadRequest)
	}
}
