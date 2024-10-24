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

// Adds new item to inventory repository
func (s *InventoryService) AddInventoryItem(item models.InventoryItem) error {
	if s.inventoryRepo.Exists(item.IngredientID) {
		return errors.New("inventory item already exists")
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

// Returns all items from inventory repository
func (s *InventoryService) GetAllInventoryItems() ([]models.InventoryItem, error) {
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		return []models.InventoryItem{}, nil
	}
	return items, nil
}

// Return item by id from inventory repository
func (s *InventoryService) GetItem(id string) (models.InventoryItem, error) {
	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}

	for _, item := range items {
		if item.IngredientID == id {
			return item, nil
		}
	}
	return models.InventoryItem{}, errors.New("inventory item does not exists")
}

// Updates Item by id replacing with new given Item in inventory repository
func (s *InventoryService) UpdateItem(id string, newItem models.InventoryItem) error {
	if !s.inventoryRepo.Exists(id) {
		return errors.New("inventory item does not exists")
	}

	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.IngredientID == id {
			items[i] = newItem
			break
		}
	}

	err = s.inventoryRepo.SaveAll(items)
	if err != nil {
		return err
	}
	return nil
}

// Deletes Item by id from inventory repository
func (s *InventoryService) DeleteItem(id string) error {
	if !s.inventoryRepo.Exists(id) {
		return errors.New("inventory item does not exists")
	}

	items, err := s.inventoryRepo.GetAll()
	if err != nil {
		return err
	}

	newItems := []models.InventoryItem{}

	for _, item := range items {
		if item.IngredientID != id {
			newItems = append(newItems, item)
		}
	}

	err = s.inventoryRepo.SaveAll(newItems)
	if err != nil {
		return err
	}

	return nil
}
