package main

import (
	"fmt"
	"hot-coffee/mainhelpers"
	"os"
)

func main() {
	Help(os.Args)
	mainhelpers.ExecuteProgram()
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
