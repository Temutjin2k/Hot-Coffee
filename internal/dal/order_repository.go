package dal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"hot-coffee/models"
)

// OrderRepository implements OrderRepository using JSON files
type OrderRepository struct {
	path string
}

// NewOrderRepository creates a new FileOrderRepository
func NewOrderRepository(path string) *OrderRepository {
	return &OrderRepository{path: path}
}

func (repo *OrderRepository) GetAll() ([]models.Order, error) {
	content, err := os.ReadFile(repo.path)
	if len(content) == 0 {
		return []models.Order{}, err
	}
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	err = json.Unmarshal(content, &orders)
	return orders, err
}

func (repo *OrderRepository) Add(order models.Order) error {
	orders, err := repo.GetAll()
	if err != nil {
		return err
	}
	orders = append(orders, order)

	jsonData, err := json.MarshalIndent(orders, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(repo.path, jsonData, os.ModePerm)
}

func (repo *OrderRepository) SaveAll(Orders []models.Order) error {
	jsonData, err := json.MarshalIndent(Orders, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.path, jsonData, 0o644)
}

func (repo *OrderRepository) GetID() (int, error) {
	configPath := filepath.Join(filepath.Dir(repo.path), "config.json")

	ConfigContent, err := os.ReadFile(configPath)
	if err != nil {
		return -1, err
	}

	var ID models.OrderID
	err = json.Unmarshal(ConfigContent, &ID)
	if err != nil {
		return -1, err
	}

	i := ID.ID
	ID.ID++
	NewContent, err := json.MarshalIndent(ID, "", "    ")
	if err != nil {
		// TODO
	}
	os.WriteFile(repo.path, NewContent, os.ModePerm)
	return i, nil
}
