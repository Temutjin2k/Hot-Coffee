package main

import (
	"fmt"
	"hot-coffee/config"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"hot-coffee/utils"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
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

	// logger init
	logFile, err := os.OpenFile(path+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(logFile, nil))

	// Initialize repositories (Data Access Layer)
	orderRepo := dal.NewOrderRepository()
	menuRepo := dal.NewMenuRepository()
	inventoryRepo := dal.NewInventoryRepository()

	// Initialize services (Business Logic Layer)
	orderService := service.NewOrderService(*orderRepo, *menuRepo)
	menuService := service.NewMenuService(*menuRepo, *inventoryRepo)
	inventoryService := service.NewInventoryService(*inventoryRepo) // TODO

	// Initialize handlers (Presentation Layer)
	orderHandler := handler.NewOrderHandler(orderService, menuService, logger)
	menuHandler := handler.NewMenuHandler(menuService, logger)
	inventoryHandler := handler.NewInventoryHandler(inventoryService, logger) // TODO
	reportHandler := handler.NewAggregationHandler(orderService, logger)

	// Setup HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", orderHandler.PostOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetOrder)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrders)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.PutOrder)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrder)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.CloseOrder)

	mux.HandleFunc("POST /menu", menuHandler.PostMenu)
	mux.HandleFunc("GET /menu", menuHandler.GetMenu)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuItem)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.PutMenuItem)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuItem)

	mux.HandleFunc("POST /inventory", inventoryHandler.PostInventory)
	mux.HandleFunc("GET /inventory", inventoryHandler.GetInventory)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetInventoryItem)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.PutInventoryItem)
	mux.HandleFunc("DELETE /inventory/{id}", inventoryHandler.DeleteInventoryItem)

	mux.HandleFunc("GET /reports/total-sales", reportHandler.TotalSalesHandler)
	mux.HandleFunc("GET /reports/popular-items", reportHandler.PopularItemsHandler)

	address := "http://localhost:" + port + "/"
	fmt.Println("Server launched on address:", address)

	logger.Info("Application started", "Address", address, "Data directory", path)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
