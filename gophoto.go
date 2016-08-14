package main

import (
	"database/sql"
	"fmt"
	"github.com/dgoodwin/gophoto/config"
	"github.com/dgoodwin/gophoto/handlers"
	"github.com/dgoodwin/gophoto/importer"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Default location of this file when run in the Docker container:
var DEFAULT_CONFIG string = "/go/src/app/gophoto-docker.yml"

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
	cfg := config.LoadConfig(configFile)

	// Establish a database connection:
	// Re-use the Goose dbconf.yml for the open string:
	db, err := sql.Open("postgres", cfg.Database.Open)
	fmt.Printf("Created db: %s\n", db)
	if err != nil {
		log.Fatal(err)
	}

	//age := 21
	//rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	/*
		var newPhotoId int
			err = db.QueryRow(`INSERT INTO media(filename, url, res_x, res_y, size)
				VALUES('something.jpg', '/something.jpg', 1200, 1024, 5019328) RETURNING id`).Scan(&newPhotoId)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created new photo: %d\n", newPhotoId)
	*/

	fmt.Printf("Hello, world.\n")
	fmt.Printf("Scanning for photos in: %s\n", cfg.ImportPath)
	var scanFunc = func(path string, f os.FileInfo, err error) error {
		checkFileErr := importer.CheckFile(path, f, db)
		return checkFileErr
	}
	err = filepath.Walk(cfg.ImportPath, scanFunc)
	fmt.Printf("Walk returned: %v\n", err)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Index)

	// If nothing else has matched attempt to serve a static file:
	// TODO: Point to a proper location for the data files? Use env var or config.
	r.PathPrefix("/").Handler(http.FileServer(
		http.Dir("/home/dev/go/src/github.com/dgoodwin/gophoto/public/")))
	log.Fatal(http.ListenAndServe(":8200", r))
}
