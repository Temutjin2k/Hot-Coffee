package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"io/ioutil"
	"net/http"
	"strings"
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

func (s *MenuService) MenuCheckByID(MenuItemID string) error {
	// Use the ProductID of the single menu item to check against existing menu items
	menuItems, _ := s.menuRepo.GetAll()
	for _, item := range menuItems {
		if item.ID == MenuItemID {
			return nil
		}
	}
	return errors.New("The requested menu item to update does not exist in menu") // Item exists in the menu // Item does not exist
}

func (s *MenuService) IngredientsCheckByID(menuItemID string, quantity int) error {
	menuItems, _ := s.menuRepo.GetAll()
	ingredientsNeeded := make(map[string]float64)
	flag := false
	for _, item := range menuItems {
		if item.ID == menuItemID {
			flag = true
			for _, ingr := range item.Ingredients {
				ingredientsNeeded[ingr.IngredientID] += float64(ingr.Quantity) * float64(quantity)
			}
		}
	}

	inventoryItems, _ := s.inventoryRepo.GetAll()

	for _, inventoryItem := range inventoryItems {
		if value, exists := ingredientsNeeded[inventoryItem.IngredientID]; exists {
			if value > inventoryItem.Quantity {
				return errors.New("Not enough ingridients for item") // Not enough ingredients
			}
		}
	}
	if flag {
		return nil
	}
	return errors.New("No ingridients for item in inventory")
}

func (s *MenuService) SubtractIngredientsByID(OrderID string, quantity int) error {
	if err := s.IngredientsCheckByID(OrderID, quantity); err != nil {
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

func (s *MenuService) CheckNewMenu(MenuItem models.MenuItem) error {
	if strings.TrimSpace(MenuItem.ID) == "" {
		return errors.New("New menu item's ID is empty")
	}
	if strings.TrimSpace(MenuItem.Name) == "" {
		return errors.New("New menu item's Name is empty")
	}
	if strings.TrimSpace(MenuItem.Description) == "" {
		return errors.New("New menu item's Description is empty")
	}
	if MenuItem.Price < 0 {
		return errors.New("New menu item's Price is awkward")
	}
	for _, ingredient := range MenuItem.Ingredients {
		if strings.TrimSpace(ingredient.IngredientID) == "" {
			return errors.New("New menu item's ingredient is empty")
		}
		if ingredient.Quantity < 0 {
			return errors.New("New menu item's quantity is awkward")
		}
	}
	return nil
}
