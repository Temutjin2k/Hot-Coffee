package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
)

type MenuService struct {
	menuRepo      dal.MenuRepository
	inventoryRepo dal.InventoryRepository
}

func NewMenuService(menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) *MenuService {
	return &MenuService{menuRepo: menuRepo, inventoryRepo: inventoryRepo}
}

func (s *MenuService) DeleteMenuItem(MenuItemID string) error {
	MenuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}
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
	s.menuRepo.SaveAll(NewMenuItems)
	return nil
}

func (s *MenuService) UpdateMenuItem(r *http.Request) error {
	RequestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var RequestedMenuItem models.MenuItem
	err = json.Unmarshal(RequestContent, &RequestedMenuItem)
	if err != nil {
		return err
	}

	MenuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return err
	}

	for i, MenuItem := range MenuItems {
		if MenuItem.ID == r.PathValue("id") {
			fmt.Println(r.PathValue("id"))
			MenuItems[i].Description = RequestedMenuItem.Description
			MenuItems[i].ID = r.PathValue("id")
			MenuItems[i].Ingredients = RequestedMenuItem.Ingredients
			MenuItems[i].Name = RequestedMenuItem.Name
			MenuItems[i].Price = RequestedMenuItem.Price
		}
	}
	s.menuRepo.SaveAll(MenuItems)
	return nil
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

func (s *MenuService) MenuCheckByID(ID string) bool {
	// Use the ProductID of the single menu item to check against existing menu items
	menuItems, _ := s.menuRepo.GetAll()
	for _, item := range menuItems {
		if item.ID == ID {
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

func (s *MenuService) IngredientsCheckByID(menuItemID string, quantity int) bool {
	menuItems, _ := s.menuRepo.GetAll()
	ingredientsNeeded := make(map[string]float64)

	for _, item := range menuItems {
		if item.ID == menuItemID {
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

func (s *MenuService) SubtractIngredientsByID(OrderID string, quantity int) error {
	if !s.IngredientsCheckByID(OrderID, quantity) {
		return errors.New("Not enough ingredients or needed ingredients do not exist")
	}

	ingredients := make(map[string]float64)
	menuItems, _ := s.menuRepo.GetAll()

	for _, item := range menuItems {
		if item.ID == OrderID {
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

func (s *MenuService) GetMenuItem(MenuItemID string) (models.MenuItem, error) {
	MenuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}
	for i, MenuItem := range MenuItems {
		if MenuItem.ID == MenuItemID {
			return MenuItems[i], nil
		}
	}
	return models.MenuItem{}, err
}

func (s *MenuService) GetMenuItems() ([]models.MenuItem, error) {
	MenuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return []models.MenuItem{}, err
	}
	return MenuItems, err
}
