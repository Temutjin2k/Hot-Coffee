package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/models"
)

// MenuRepository implements MenuRepository using JSON files
type MenuRepository struct {
	path string
}

// NewMenuRepository creates a new FileMenuRepository
func NewMenuRepository(path string) *MenuRepository {
	return &MenuRepository{path: path}
}

func (repo *MenuRepository) GetAll() ([]models.MenuItem, error) {
	content, err := os.ReadFile(repo.path)
	if err != nil {
		return nil, err
	}

	var menuItems []models.MenuItem
	err = json.Unmarshal(content, &menuItems)
	return menuItems, err
}

func (repo *MenuRepository) Exists(itemID string) bool {
	items, _ := repo.GetAll()
	for _, item := range items {
		if item.ID == itemID {
			return true
		}
	}
	return false
}

func (repo *MenuRepository) SaveAll(menuItems []models.MenuItem) error {
	jsonData, err := json.MarshalIndent(menuItems, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.path, jsonData, 0o644)
}
