package MenuHandlers

import (
	"encoding/json"
	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
)

func MenuDelete(w http.ResponseWriter, MenuItemID string) {
	MenuContents, err := ioutil.ReadFile(config.BaseDir + "/menu_items.json")
	if err != nil {
		ErrorHandler.Error(w, "Could not read menu items from server", http.StatusInternalServerError)
		return
	}
	var MenuItems []models.MenuItem
	json.Unmarshal(MenuContents, &MenuItems)

	NewMenuItems := make([]models.MenuItem, 0)

	for _, MenuItem := range MenuItems {
		if MenuItem.ID != MenuItemID {
			var NewItem models.MenuItem
			NewItem.Description = MenuItem.Description
			NewItem.ID = MenuItem.ID
			NewItem.Ingredients = MenuItem.Ingredients
			NewItem.Name = MenuItem.Name
			NewItem.Price = MenuItem.Price
			NewMenuItems = append(NewMenuItems, NewItem)
		}
	}
	jsonData, err := json.MarshalIndent(NewMenuItems, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could not transfer menu items to json file", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile(config.BaseDir+"/menu_items.json", jsonData, os.ModePerm)
	if err != nil {
		ErrorHandler.Error(w, "Could not rewwrite menu items in server", http.StatusInternalServerError)
		return
	}
}
