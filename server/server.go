package server

import (
	"database/sql"
	"fmt"
	"github.com/dgoodwin/gophoto/config"
	"github.com/dgoodwin/gophoto/server/handlers"
	"github.com/dgoodwin/gophoto/server/handlers/api"
	"github.com/dgoodwin/gophoto/server/importer"
	"github.com/dgoodwin/gophoto/server/storage"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
)

// Default location of this file when run in the Docker container:
var DEFAULT_CONFIG string = "/go/src/app/gophoto-docker.yml"

func RunServer(cmd *cobra.Command, args []string) {

	log.SetLevel(log.DebugLevel)
	log.Infoln("Launching GoPhoto server")

	// Load the config file which we assume to be the only argument given on the CLI,
	// or the default file location in the Docker container if it exists.
	var configFile string = ""
	if len(args) == 1 {
		configFile, _ = filepath.Abs(args[0])
	} else if _, err := os.Stat(DEFAULT_CONFIG); len(os.Args) < 2 && err == nil {
		configFile = DEFAULT_CONFIG
	} else {
		fmt.Printf("No config file given on CLI or in default location: %s\n", DEFAULT_CONFIG)
		os.Exit(1)
	}

	log.Infof("Loading config from: %s\n", configFile)
	cfg := config.LoadConfigFile(configFile)

	// Establish a database connection:
	// Re-use the Goose dbconf.yml for the open string:
	db, err := sql.Open("postgres", cfg.Database.Open)
	log.Debugf("Created db connection: %s", db)
	if err != nil {
		log.Fatal(err)
	}

	var s storage.StorageBackend
	s, err = storage.NewStorageBackend(cfg)

	i := importer.Importer{DB: db, Storage: s}

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Index)

	r.Handle("/api/v1/media", api.MediaHandler{Cfg: cfg, Importer: i}).Methods("POST")

	// If nothing else has matched attempt to serve a static file:
	// TODO: Point to a proper location for the data files? Use env var or config.
	r.PathPrefix("/").Handler(http.FileServer(
		http.Dir("/home/dev/go/src/github.com/dgoodwin/gophoto/public/")))

	log.Fatal(http.ListenAndServe(":8200", r))
	log.Infoln("Listening on port 8200")
}
