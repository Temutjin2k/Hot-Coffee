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
	// Check for help flag
	utils.Help(os.Args)

	user, err := user.Current() // Создаем объект user, а уже из него достаем основную директорию user.HomeDir
	if err != nil {
		fmt.Println("Error getting user home directory")
		os.Exit(1)
	}

	// Checking Flags
	dir, port := config.Flagchecker()
	path := filepath.Join(user.HomeDir, "hot-coffee", dir)

	if !utils.DirectoryExists(path) {
		utils.CreateDir(dir)
	}

	// Initialize repositories (Data Access Layer)
	menuRepo := dal.NewMenuRepository()
	inventoryRepo := dal.NewInventoryRepository()
	orderRepo := dal.NewOrderRepository()

	// Initialize services (Business Logic Layer)
	menuService := service.NewMenuService(*menuRepo, *inventoryRepo)
	orderService := service.NewOrderService(*orderRepo)
	inventoryService := service.NewInventoryService(*inventoryRepo) // TODO

	// Initialize handlers (Presentation Layer)
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService) // TODO
	reportHandler := handler.NewAggregationHandler(orderService)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/menu", menuHandler.MenuHandler)
	mux.HandleFunc("/orders", orderHandler.OrderHandler)
	mux.HandleFunc("/orders/{id}", orderHandler.OrderHandler)
	mux.HandleFunc("/inventory", inventoryHandler.InventoryHandler)
	mux.HandleFunc("/reports/total-sales", reportHandler.TotalSalesHandler)
	mux.HandleFunc("/reports/popular-items", reportHandler.PopularItemsHandler)

	address := "http://localhost:" + port + "/"
	fmt.Println("Server launched on address:", address)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
