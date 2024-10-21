package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"hot-coffee/config"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"hot-coffee/utils"
)

func main() {
	utils.Help(os.Args)
	user, err := user.Current() // Создаем объект user, а уже из него достаем основную директорию user.HomeDir
	if err != nil {
		fmt.Println("Error getting user home directory")
		os.Exit(1)
	}

	dir, port := config.Flagchecker()
	path := filepath.Join(user.HomeDir, "hot-coffee", dir)

	if !utils.DirectoryExists(path) {
		utils.CreateDir(dir)
	}
	// Initialize repositories
	menuRepo := dal.NewMenuRepository()
	inventoryRepo := dal.NewInventoryRepository()
	orderRepo := dal.NewOrderRepository() // Initialize the order repository

	// Initialize services
	menuService := service.NewMenuService(*menuRepo, *inventoryRepo)
	orderService := service.NewOrderService(*orderRepo) // Initialize the order service

	// Initialize handlers
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService) // Pass the order service
	inventoryHandler := handler.NewInventoryHandler()
	reportsHandler := handler.NewReportsHandler()

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/menu", menuHandler.MenuHandler)
	mux.HandleFunc("/orders", orderHandler.OrderHandler)
	mux.HandleFunc("/orders/", orderHandler.OrderHandler)
	mux.HandleFunc("/inventory", inventoryHandler.InventoryHandler)
	mux.HandleFunc("/reports/total-sales", handler.TotalSalesHandler)
	mux.HandleFunc("/reports/popular-items", handler.PopularItemsHandler)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
