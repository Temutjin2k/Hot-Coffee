package services

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/ErrorHandler"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"os"
)

func MenuCheck(w http.ResponseWriter, order models.Order) bool { // Надо проверить если в меню эта вещь
	items := make([]string, 0)
	for _, item := range order.Items {
		items = append(items, item.ProductID)
	}
	fmt.Println(items) // espresso

	var MenuItems []models.MenuItem
	menucontent, err := ioutil.ReadFile("data/menu_items.json")
	if err != nil {
		// TO DO
	}

	json.Unmarshal(menucontent, &MenuItems)
	match := 0
	for i := 0; i < len(items); i++ {
		for _, item := range MenuItems {
			if item.ID == items[i] {
				match++
			}
		}
	}
	return match == len(items)
}

func IngredientsCheck(w http.ResponseWriter, order models.Order) bool { // Проверка на ингредиенты
	if !MenuCheck(w, order) {
		ErrorHandler.Error(w, "Your order is not in menu, please check again our menu", http.StatusBadRequest)
		return false
	}
	menucontent, err := ioutil.ReadFile("data/menu_items.json")
	if err != nil {
		// TO DO
	}

	var MenuItems []models.MenuItem
	json.Unmarshal(menucontent, &MenuItems)
	ing := make(map[string]float64)

	for _, orderitem := range order.Items {
		for _, menuitem := range MenuItems {
			if orderitem.ProductID == menuitem.ID {
				for _, ingrs := range menuitem.Ingredients {
					ing[ingrs.IngredientID] = float64(ingrs.Quantity) * float64(orderitem.Quantity)
				}
			}
		}
	}

	inventorycontent, err := ioutil.ReadFile("data/inventory.json")
	if err != nil {
		// TO DO
	}
	fmt.Println(ing)
	flag := true
	var InventoryItems []models.InventoryItem
	json.Unmarshal(inventorycontent, &InventoryItems)

	for _, inventoryitem := range InventoryItems {
		value, isExist := ing[inventoryitem.IngredientID]
		if isExist {
			if value < inventoryitem.Quantity {
			} else {
				flag = false
			}
		}
	}

	return flag
}

func SubtractIngridients(w http.ResponseWriter, order models.Order) {
	if !IngredientsCheck(w, order) {
		ErrorHandler.Error(w, "Not enough ingedients or needed ingerdients do not exist", http.StatusBadRequest)
		return
	}
	menucontent, err := ioutil.ReadFile("data/menu_items.json")
	if err != nil {
		// TO DO
	}

	var MenuItems []models.MenuItem
	json.Unmarshal(menucontent, &MenuItems)
	ing := make(map[string]float64)

	for _, orderitem := range order.Items {
		for _, menuitem := range MenuItems {
			if orderitem.ProductID == menuitem.ID {
				for _, ingrs := range menuitem.Ingredients {
					ing[ingrs.IngredientID] = float64(ingrs.Quantity) * float64(orderitem.Quantity)
				}
			}
		}
	}

	inventorycontent, err := ioutil.ReadFile("data/inventory.json")
	if err != nil {
		// TO DO
	}
	fmt.Println(ing)

	var InventoryItems []models.InventoryItem
	json.Unmarshal(inventorycontent, &InventoryItems)

	for i, inventoryitem := range InventoryItems {
		value, isExist := ing[inventoryitem.IngredientID]
		if isExist {
			if value < inventoryitem.Quantity {
				InventoryItems[i].Quantity -= value
			}
		}
	}

	// json.MarshalIndent() принимает структуру и конвертирует все в инфу в стиле json
	jsonData, err := json.MarshalIndent(InventoryItems, "", "    ")
	if err != nil {
		// will write error handler
	}

	err = os.WriteFile("data/inventory.json", jsonData, 0644) // os.WriteFile(filename, content, perm) в файл записывает данные
	if err != nil {
		// will write error handler
	}
}
