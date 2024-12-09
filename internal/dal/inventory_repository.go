package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

type InventoryRepo interface {
	GetAll() ([]models.InventoryItem, error)
	Exists(ID string) bool
	SubtractIngredients(ingredients map[string]float64) error
	SaveAll(items []models.InventoryItem) error
}

// InventoryRepository implements InventoryRepository using JSON files
type InventoryRepository struct {
	path string
}

// NewInventoryRepository creates a new FileInventoryRepository
func NewInventoryRepository(path string) *InventoryRepository {
	return &InventoryRepository{path: path}
}

func (repo *InventoryRepository) GetAll() ([]models.InventoryItem, error) {
	content, err := os.ReadFile(repo.path)
	if err != nil {
		return nil, err
	}

	var inventoryItems []models.InventoryItem
	err = json.Unmarshal(content, &inventoryItems)
	return inventoryItems, err
}

// Check if ingridient by given ID exists
func (repo *InventoryRepository) Exists(ID string) bool {
	inventoryItems, err := repo.GetAll()
	if err != nil {
		return false
	}

	for _, item := range inventoryItems {
		if item.IngredientID == ID {
			return true
		}
	}
	return false
}

func (repo *InventoryRepository) SubtractIngredients(ingredients map[string]float64) error {
	inventoryItems, err := repo.GetAll()
	if err != nil {
		return err
	}

	for i, inventoryItem := range inventoryItems {
		if value, exists := ingredients[inventoryItem.IngredientID]; exists {
			inventoryItems[i].Quantity -= value
		}
	}

	jsonData, err := json.MarshalIndent(inventoryItems, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(repo.path, jsonData, 0o644)
}

func (repo *InventoryRepository) SaveAll(items []models.InventoryItem) error {
	jsonData, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.path, jsonData, 0o644)
}
