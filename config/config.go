package config

import "flag"

var BaseDir string

func Flagchecker() (string, string) {
	dir := flag.String("dir", "data", "Path to the data directory")
	port := flag.String("port", "8080", "Port Number")
	BaseDir = *dir
	flag.Parse()
	return *dir, *port
}
