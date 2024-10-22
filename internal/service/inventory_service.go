package service

import (
	"errors"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

// Some funcs
// ...

func (s *InventoryService) AddInventoryItem(item models.InventoryItem) error {
	if s.inventoryRepo.Exists(item.IngredientID) {
		return errors.New("inventory item, already exists")
	}

	ingridientItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}

	ingridientItems = append(ingridientItems, item)

	err = s.inventoryRepo.SaveAll(ingridientItems)
	if err != nil {
		return err
	}

	return nil
}
