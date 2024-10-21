package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/config"
	"hot-coffee/models"
)

// MenuRepository implements MenuRepository using JSON files
type MenuRepository struct{}

// NewMenuRepository creates a new FileMenuRepository
func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

func (repo *MenuRepository) GetAll() ([]models.MenuItem, error) {
	content, err := os.ReadFile(config.BaseDir + "/menu_items.json")
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
	return os.WriteFile(config.BaseDir+"/menu_items.json", jsonData, 0o644)
}
