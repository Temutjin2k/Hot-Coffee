package main

import (
	"flag"
	"fmt"
	"hot-coffee/internal/handler"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	Help(os.Args)

	dir := flag.String("dir", "maks", "Path to the data directory")
	port := flag.String("port", "8080", "Port Number")
	flag.Parse()

	CreateDir(*dir)

	http.HandleFunc("/orders", handler.OrderHandler)
	http.HandleFunc("/menu", handler.MenuHandler)
	http.HandleFunc("/inventory", handler.InventoryHandler)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Println("Error starting server")
		return
	}
}

var BaseDir string

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
		fmt.Println("Something went wrong when creating a data directory")
		os.Exit(1)
	}
}

func Help(s []string) {
	for _, v := range s {
		if v == "--help" || v == "-help" || v == "-h" {
			fmt.Print(`Coffee Shop Management System

Usage:
	hot-coffee [--port <N>] [--dir <S>] 
	hot-coffee --help
			
Options:
	--help       Show this screen.
	--port N     Port number.
	--dir S      Path to the data directory.`)
		}
	}
}
