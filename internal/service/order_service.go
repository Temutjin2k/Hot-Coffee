package service

import (
	"errors"

	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type OrderService struct {
	orderRepo dal.OrderRepository
}

func NewOrderService(orderRepo dal.OrderRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo}
}

// AddOrder adds a new order to the repository
func (s *OrderService) AddOrder(order models.Order) error {
	if order.ID == "" {
		return errors.New("order ID cannot be empty")
	}

	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	// Check for duplicate order ID
	for _, existingOrder := range existingOrders {
		if existingOrder.ID == order.ID {
			return errors.New("order with this ID already exists")
		}
	}

	return s.orderRepo.Add(order)
}

// GetAllOrders retrieves all orders from the repository
func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAll()
}

// DeleteOrder removes an order by ID
func (s *OrderService) DeleteOrder(orderID string) error {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for _, order := range existingOrders {
		if order.ID == orderID {
			return s.orderRepo.Delete(orderID)
		}
	}

	return errors.New("order not found")
}

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(updatedOrder models.Order) error {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for _, order := range existingOrders {
		if order.ID == updatedOrder.ID {
			return s.orderRepo.Update(updatedOrder)
		}
	}

	return errors.New("order not found")
}

func (s *OrderService) GetTotalSales() (models.TotalSales, error) {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return models.TotalSales{}, err
	}

	// Counting sales amount
	totalSales := models.TotalSales{}

	for _, order := range existingOrders {
		for _, item := range order.Items {
			totalSales.TotalSales += item.Quantity
		}
	}
	return totalSales, nil
}
