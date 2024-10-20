package MenuHandlers

import (
	"encoding/json"
	"fmt"
	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/services"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
)

func MenuPut(w http.ResponseWriter, r *http.Request, MenuItemID string) {
	RequestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var RequestedMenuItem models.MenuItem
	err = json.Unmarshal(RequestContent, &RequestedMenuItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	if !services.MenuCheck2(RequestedMenuItem) {
		ErrorHandler.Error(w, "The menu item you want to update does not exist in menu", http.StatusNotFound)
	}

	MenuContent, err := ioutil.ReadFile(config.BaseDir + "/menu_items.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read menu items from server", http.StatusInternalServerError)
		return
	}
	var MenuItems []models.MenuItem
	json.Unmarshal(MenuContent, &MenuItems)

	for i, MenuItem := range MenuItems {
		if MenuItem.ID == MenuItemID {
			fmt.Println(MenuItemID)
			MenuItems[i].Description = RequestedMenuItem.Description
			MenuItems[i].ID = MenuItemID
			MenuItems[i].Ingredients = RequestedMenuItem.Ingredients
			MenuItems[i].Name = RequestedMenuItem.Name
			MenuItems[i].Price = RequestedMenuItem.Price
		}
	}

	jsondata, err := json.MarshalIndent(MenuItems, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could not upload menu items to server", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile(config.BaseDir+"/menu_items.json", jsondata, os.ModePerm)
	if err != nil {
		ErrorHandler.Error(w, "Could not rewrite orders", http.StatusInternalServerError)
		return
	}
}
