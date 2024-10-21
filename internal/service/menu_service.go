package service

import (
	"errors"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type MenuService struct {
	menuRepo      dal.MenuRepository
	inventoryRepo dal.InventoryRepository
}

func NewMenuService(menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) *MenuService {
	return &MenuService{menuRepo: menuRepo, inventoryRepo: inventoryRepo}
}

func (s *MenuService) MenuCheck(menuItem models.MenuItem) bool {
	// Use the ProductID of the single menu item to check against existing menu items
	menuItems, _ := s.menuRepo.GetAll()
	for _, item := range menuItems {
		if item.ID == menuItem.ID {
			return true // Item exists in the menu
		}
	}
	return false // Item does not exist
}

func (s *MenuService) IngredientsCheck(menuItem models.MenuItem, quantity int) bool {
	menuItems, _ := s.menuRepo.GetAll()
	ingredientsNeeded := make(map[string]float64)

	for _, item := range menuItems {
		if item.ID == menuItem.ID {
			for _, ingr := range item.Ingredients {
				ingredientsNeeded[ingr.IngredientID] += float64(ingr.Quantity) * float64(quantity)
			}
		}
	}

	inventoryItems, _ := s.inventoryRepo.GetAll()

	for _, inventoryItem := range inventoryItems {
		if value, exists := ingredientsNeeded[inventoryItem.IngredientID]; exists {
			if value > inventoryItem.Quantity {
				return false // Not enough ingredients
			}
		}
	}

	return true // Enough ingredients available
}

func (s *MenuService) SubtractIngredients(menuItem models.MenuItem, quantity int) error {
	if !s.IngredientsCheck(menuItem, quantity) {
		return errors.New("not enough ingredients or needed ingredients do not exist")
	}

	ingredients := make(map[string]float64)
	menuItems, _ := s.menuRepo.GetAll()

	for _, item := range menuItems {
		if item.ID == menuItem.ID {
			for _, ingr := range item.Ingredients {
				ingredients[ingr.IngredientID] += float64(ingr.Quantity) * float64(quantity)
			}
		}
	}

	return s.inventoryRepo.SubtractIngredients(ingredients)
}

func (s *MenuService) AddMenuItem(menuItem models.MenuItem) error {
	// Check if the item already exists
	if s.MenuCheck(menuItem) {
		return errors.New("menu item already exists")
	}

	// Load current menu items
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}

	// Append the new menu item to the existing list
	menuItems = append(menuItems, menuItem)

	// Save the updated list back to the repository
	return s.menuRepo.SaveAll(menuItems)
}
