package menuHandler

import (
	"encoding/json"
	"net/http"
	"os"

	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/models"
)

func GetMenuItems(w http.ResponseWriter) {
	MenuContent, err := os.ReadFile(config.BaseDir + "/menu_items.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read menu items", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(MenuContent)
}

func GetMenuItem(w http.ResponseWriter, MenuItemID string) {
	MenuContent, err := os.ReadFile(config.BaseDir + "/menu_items.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read menu items", http.StatusInternalServerError)
	}
	var MenuItems []models.MenuItem
	err = json.Unmarshal(MenuContent, &MenuItems)
	if err != nil {
		ErrorHandler.Error(w, "Could not convert menu data to json", http.StatusInternalServerError)
	}
	flag := true
	var NeededMenuItem models.MenuItem
	for i, MenuItem := range MenuItems {
		if MenuItem.ID == MenuItemID {
			flag = false
			NeededMenuItem = MenuItems[i]
		}
	}
	jsondata, err := json.Marshal(NeededMenuItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not convert menu data to json", http.StatusInternalServerError)
	}
	if !flag {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsondata)
	} else {
		ErrorHandler.Error(w, "Your requested menu item is not found", 404)
	}
}
