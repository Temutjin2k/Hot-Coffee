package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	Parts := strings.SplitN(r.URL.Path[1:], "/", 3)
	switch len(Parts) {
	case 1:
		switch r.Method {
		case http.MethodPost:
			PostOrder(w, r)
		case http.MethodGet:
			GetOrders(w, r)
		default:
			fmt.Println("error")
			os.Exit(1)

		}
	case 2:
		switch r.Method {
		case http.MethodPost:

		case http.MethodGet:

		case http.MethodDelete:

		default:
			fmt.Println("error")
			os.Exit(1)
		}
	case 3:
		if r.Method == http.MethodPost {
		} else {
			fmt.Println("error")
			os.Exit(1)
		}
	}
}
