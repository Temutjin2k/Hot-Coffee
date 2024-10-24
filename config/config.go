package config

import "hot-coffee/models"

var (
	BaseDir string
	// Default inventory data
	DefaultInventoryData = []models.InventoryItem{
		{IngredientID: "espresso_shot", Name: "Espresso Shot", Quantity: 500, Unit: "shots"},
		{IngredientID: "milk", Name: "Milk", Quantity: 5000, Unit: "ml"},
		{IngredientID: "flour", Name: "Flour", Quantity: 6000, Unit: "g"},
		{IngredientID: "blueberries", Name: "Blueberries", Quantity: 2000, Unit: "g"},
		{IngredientID: "sugar", Name: "Sugar", Quantity: 5000, Unit: "g"},
	}

	// Default menu items data
	DefaultMenuItemsData = []models.MenuItem{
		{
			ID:          "latte",
			Name:        "Caffe Latte",
			Description: "Espresso with steamed milk",
			Price:       3.5,
			Ingredients: []models.MenuItemIngredient{
				{IngredientID: "espresso_shot", Quantity: 1},
				{IngredientID: "milk", Quantity: 200},
			},
		},
		{
			ID:          "muffin",
			Name:        "Blueberry Muffin",
			Description: "Freshly baked muffin with blueberries",
			Price:       2,
			Ingredients: []models.MenuItemIngredient{
				{IngredientID: "flour", Quantity: 100},
				{IngredientID: "blueberries", Quantity: 20},
				{IngredientID: "sugar", Quantity: 30},
			},
		},
		{
			ID:          "espresso",
			Name:        "Espresso",
			Description: "Strong and bold coffee",
			Price:       2.5,
			Ingredients: []models.MenuItemIngredient{
				{IngredientID: "espresso_shot", Quantity: 1},
			},
		},
	}

	DefaultOrdersData = []models.MenuItem{}

	DefaultConfigData = models.OrderID{ID: 0}
)
