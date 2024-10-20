package MenuHandlers

import (
	"encoding/json"
	"hot-coffee/config"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/internal/services"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
)

func MenuPost(w http.ResponseWriter, r *http.Request) {
	var NewItem models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&NewItem)
	if err != nil {
		ErrorHandler.Error(w, "Could not decode request json data", http.StatusBadRequest)
		return
	}

	if services.MenuCheck2(NewItem) {
		ErrorHandler.Error(w, "The requested menu item already exists in current menu", http.StatusBadRequest)
		return
	}

	var MenuItems []models.MenuItem
	content, err := ioutil.ReadFile(config.BaseDir + "/menu_items.json") // ioutil.ReadFile() читает файл и возвращает содержимое в массие из байтов
	if err != nil {
		ErrorHandler.Error(w, "Could not read orders from server", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(content, &MenuItems) // json.Unmarshal([]byte, type any) короче из инфы в байтах он конвертирует все в структуру в стиле json

	MenuItems = append(MenuItems, NewItem)

	// json.MarshalIndent() принимает структуру и конвертирует все в инфу в стиле json
	jsonData, err := json.MarshalIndent(MenuItems, "", "    ")
	if err != nil {
		ErrorHandler.Error(w, "Could not convert orders to json data", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(config.BaseDir+"/menu_items.json", jsonData, 0644) // os.WriteFile(filename, content, perm) в файл записывает данные
	if err != nil {
		ErrorHandler.Error(w, "Could not write orders to json database", http.StatusInternalServerError)
		return
	}
}
