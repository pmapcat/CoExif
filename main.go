package main

import (
	"flag"
)

var BASE_SIZE = 10
var AUTO_SPAWN = true
var EXIF_DIR = "./exif_tool/exiftool"
var TEMP_DIR = "./temp"

func main() {
	port := flag.String("port", "9999", "Enter a server port number")
	root_dir := flag.String("root", "/", "Enter default root path")
	auth_name := flag.String("auth-name", "admin", "Enter auth name")
	auth_pass := flag.String("auth-pass", "admin", "Enter auth pass")
	max_procs := flag.Int("max-prox", 10, "Enter number of ExifTool processes")
	exif_dir := flag.String("exif-path", "./exif_tool/exiftool", "Enter path to exiftool")
	auto_spawn := flag.Bool("auto-spawn", false, "Should I autospawn processes")
	flag.Parse()

	BASE_SIZE = *max_procs
	AUTO_SPAWN = *auto_spawn
	EXIF_DIR = *exif_dir
	r := RESTHandler{}
	r.Run(*port, *auth_name, *auth_pass, *root_dir)
}
