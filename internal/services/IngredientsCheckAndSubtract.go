package services

import (
	"encoding/json"
	"hot-coffee/models"
	"io/ioutil"
)

func MenuCheck(order models.Order) bool { // Надо проверить если в меню эта вещь
	items := make([]string, 0)
	for _, item := range order.Items {
		items = append(items, item.ProductID)
	}
	var MenuItems []models.MenuItem
	menucontent, err := ioutil.ReadFile("data/menu-items.json")
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

func IngredientsCheck(order models.Order) bool { // Проверка на ингредиенты
	if !MenuCheck(order) {
		// вывести ошибку что нету этой вещи в меню
	}
	ingredients, err := ioutil.ReadFile("data/inventory.json")
	if err != nil {
		// TO DO
	}
	var MenuItems []models.MenuItem
	json.Unmarshal(ingredients, &MenuItems)

	return true
}

// func SubtractIngredient() { // Отнять ингредиенты
// }
