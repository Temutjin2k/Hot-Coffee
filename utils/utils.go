package utils

import (
	"fmt"
	"hot-coffee/config"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func Help(s []string) {
	for _, v := range s {
		if v == "--help" || v == "-help" || v == "-h" {
			fmt.Println(`Coffee Shop Management System

Usage:
	hot-coffee [--port <N>] [--dir <S>] 
	hot-coffee --help
			
Options:
	--help       Show this screen.
	--port N     Port number.
	--dir S      Path to the data directory.`)
			os.Exit(0)
		}
	}
}

func CreateDir(dir string) {
	config.BaseDir = dir

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
		dstPath := filepath.Join(user.HomeDir, "hot-coffee", config.BaseDir, entry.Name())
		err := copyFile(scrPath, dstPath)
		if err != nil {
			fmt.Println("Something went wrong when creating a data directory")
			os.Exit(1)
		}
	}
	scrPath := filepath.Join(user.HomeDir, "hot-coffee", "config", "config.json")
	dstPath := filepath.Join(user.HomeDir, "hot-coffee", config.BaseDir, "config.json")
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
