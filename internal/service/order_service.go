package service

import (
	"errors"
	"sort"

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

// Returns Popular Items sorted in decreasing order. Number of returned items depends on passing value(popularItemsNum)
func (s *OrderService) GetPopularItems(popularItemsNum int) (models.PopularItems, error) {
	existingOrders, err := s.orderRepo.GetAll()
	if err != nil {
		return models.PopularItems{}, err
	}

	// Should return sorted decreasing array
	itemMap := make(map[string]int)
	for _, order := range existingOrders {
		for _, item := range order.Items {
			itemMap[item.ProductID] += item.Quantity
		}
	}

	sortedItems := make([]models.OrderItem, 0, len(itemMap))
	for productID, quantity := range itemMap {
		sortedItems = append(sortedItems, models.OrderItem{ProductID: productID, Quantity: quantity})
	}

	// Sorting in decresing order
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].Quantity > sortedItems[j].Quantity
	})

	// To prevent from out of range
	if popularItemsNum > len(sortedItems) {
		popularItemsNum = len(sortedItems)
	}

	popularItems := models.PopularItems{Items: sortedItems[:popularItemsNum]} // potential out of range
	return popularItems, nil
}
