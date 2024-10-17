package businesslogic

import (
	"hot-coffee/models"
	"io/ioutil"
)

func MenuCheck(order models.Order) bool { // Надо проверить если в меню эта вещб
}

func IngredientsCheck(order models.Order) bool { // Проверка на ингредиенты
	ingredients, err := ioutil.ReadFile("data/inventory.json")
	if err != nil {
		// TO DO
	}
	return true
}

func SubIngredient() { // Отнять ингредиенты
}
