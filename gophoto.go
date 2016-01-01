package main

import (
	"fmt"
	"github.com/dgoodwin/gophoto/config"
	"github.com/dgoodwin/gophoto/handlers"
	"github.com/gorilla/mux"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func checkFile(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && isImage(path) {
		fmt.Printf("Visited: %s\n", path)
		//fmt.Printf("  %d %s %s\n", f.Size(), f.Mode(), f.IsDir())
		width, height := getImageDimension(path)
		fmt.Println("  Width:", width, "  Height:", height)
	}
	return nil
}

func isImage(path string) bool {
	imageExtensions := map[string]bool{
		".jpg": true,
		".JPG": true,
	}
	if imageExtensions[filepath.Ext(path)] {
		return true
	}
	return false
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

func main() {

	// Load the config file which we assume to be the only argument given on the CLI.
	if len(os.Args) != 2 {
		fmt.Printf("Must specify config a gophoto config file on CLI.\n")
		os.Exit(1)
	}
	configFile, _ := filepath.Abs(os.Args[1])
	fmt.Printf("Loading config from: %s\n", configFile)
	config.LoadConfig(configFile)

	fmt.Printf("Hello, world.\n")
	importDir := "/home/dev/Photos/2015/11"
	err := filepath.Walk(importDir, checkFile)
	fmt.Printf("Walk returned: %v\n", err)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Index)

	// If nothing else has matched attempt to serve a static file:
	// TODO: Point to a proper location for the data files? Use env var or config.
	r.PathPrefix("/").Handler(http.FileServer(
		http.Dir("/home/dev/go/src/github.com/dgoodwin/gophoto/public/")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
