package Execution

import (
	"fmt"
	"hot-coffee/config"
	"hot-coffee/internal/Mainhandlers"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

var BaseDir string

func ExecuteProgram() {
	user, err := user.Current() // Создаем объект user, а уже из него достаем основную директорию user.HomeDir
	if err != nil {
		fmt.Println("Error getting user home directory")
		os.Exit(1)
	}
	dir, port := config.Flagchecker()
	path := filepath.Join(user.HomeDir, "hot-coffee", dir)
	if !directoryExists(path) {
		CreateDir(dir)
	}
	http.HandleFunc("/menu", Mainhandlers.MenuHandler)
	http.HandleFunc("/orders", Mainhandlers.OrderHandler)
	http.HandleFunc("/orders/", Mainhandlers.OrderHandler)
	http.HandleFunc("/inventory", Mainhandlers.InventoryHandler)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server")
		return
	}
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func CreateDir(dir string) {
	BaseDir = dir

	user, err := user.Current() // Создаем объект user, а уже из него достаем основную директорию user.HomeDir
	if err != nil {
		fmt.Println("Error getting user home directory")
		os.Exit(1)
	}

	path := filepath.Join(user.HomeDir, "hot-coffee", dir)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println("Something went wrong when 5creating a data directory")
		os.Exit(1)
	}

	path = filepath.Join(user.HomeDir, "hot-coffee", "data")
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Something went wrong when creating a data directory")
		os.Exit(1)
	}

	for _, entry := range entries {
		scrPath := filepath.Join(user.HomeDir, "hot-coffee", "data", entry.Name())
		dstPath := filepath.Join(user.HomeDir, "hot-coffee", BaseDir, entry.Name())
		err := copyFile(scrPath, dstPath)
		if err != nil {
			fmt.Println("Something went wrong when creating a data directory")
			os.Exit(1)
		}
	}
	scrPath := filepath.Join(user.HomeDir, "hot-coffee", "config", "config.json")
	dstPath := filepath.Join(user.HomeDir, "hot-coffee", BaseDir, "config.json")
	copyFile(scrPath, dstPath)
}

func copyFile(src string, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		fmt.Println("Something went wrong when creating a data directory")
		os.Exit(1)
	}
	defer sourceFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("Something went wrong when creating a data directory")
		os.Exit(1)
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}
