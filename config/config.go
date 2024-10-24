package config

import "flag"

var BaseDir string

func Flagchecker() (string, string) {
	dir := flag.String("dir", "data", "Path to the data directory")
	port := flag.String("port", "8080", "Port Number")
	flag.Parse()
	BaseDir = *dir
	return *dir, *port
}
