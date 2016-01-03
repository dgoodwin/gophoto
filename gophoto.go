package main

import (
	"fmt"
	"github.com/dgoodwin/gophoto/config"
	"github.com/dgoodwin/gophoto/handlers"
	"github.com/dgoodwin/gophoto/importer"
	"github.com/gorilla/mux"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Default location of this file when run in the Docker container:
var DEFAULT_CONFIG string = "/go/src/app/devel-config.yml"

func main() {

	// Load the config file which we assume to be the only argument given on the CLI,
	// or the default file location in the Docker container if it exists.
	var configFile string = ""
	if len(os.Args) == 2 {
		configFile, _ = filepath.Abs(os.Args[1])
	} else if _, err := os.Stat(DEFAULT_CONFIG); len(os.Args) < 2 && err == nil {
		configFile = DEFAULT_CONFIG
	} else {
		fmt.Printf("No config file given on CLI or in default location: %s", DEFAULT_CONFIG)
		os.Exit(1)
	}

	fmt.Printf("Loading config from: %s\n", configFile)
	config.LoadConfig(configFile)

	fmt.Printf("Hello, world.\n")
	importDir := "/import"
	err := filepath.Walk(importDir, importer.CheckFile)
	fmt.Printf("Walk returned: %v\n", err)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Index)

	// If nothing else has matched attempt to serve a static file:
	// TODO: Point to a proper location for the data files? Use env var or config.
	r.PathPrefix("/").Handler(http.FileServer(
		http.Dir("/home/dev/go/src/github.com/dgoodwin/gophoto/public/")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
