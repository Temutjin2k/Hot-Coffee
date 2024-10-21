package dal

import (
	"encoding/json"
	"os"

	"hot-coffee/config"
	"hot-coffee/models"
)

// OrderRepository implements OrderRepository using JSON files
type OrderRepository struct{}

// NewOrderRepository creates a new FileOrderRepository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (repo *OrderRepository) GetAll() ([]models.Order, error) {
	content, err := os.ReadFile(config.BaseDir + "/orders.json")
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

	return os.WriteFile(config.BaseDir+"/orders.json", jsonData, 0o644)
}

func (repo *OrderRepository) Delete(orderID string) error {
	return nil
}

func (repo *OrderRepository) Update(order models.Order) error {
	return nil
}
