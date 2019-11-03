package server

import (
	"database/sql"
	"fmt"
	"github.com/dgoodwin/gophoto/config"
	"github.com/dgoodwin/gophoto/pkg/api/v1/dbclient"
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

// defaultConfig is the location of the config file when run in the Docker container:
var defaultConfig = "/go/src/app/gophoto-docker.yml"

// RunServer starts the gophoto server components.
func RunServer(cmd *cobra.Command, args []string) {

	log.SetLevel(log.DebugLevel)
	log.Infoln("Launching GoPhoto server")

	// Load the config file which we assume to be the only argument given on the CLI,
	// or the default file location in the Docker container if it exists.
	var configFile string
	if len(args) == 1 {
		configFile, _ = filepath.Abs(args[0])
	} else if _, err := os.Stat(defaultConfig); len(os.Args) < 2 && err == nil {
		configFile = defaultConfig
	} else {
		fmt.Printf("No config file given on CLI or in default location: %s\n", defaultConfig)
		os.Exit(1)
	}

	log.Infof("Loading config from: %s", configFile)
	cfg := config.LoadConfigFile(configFile)
	log.Infoln(cfg)

	srv := buildServer(cfg)
	log.Fatal(srv.ListenAndServe())
}

func buildServer(cfg config.GophotoConfig) *http.Server {
	// Establish a database connection:
	// Re-use the Goose dbconf.yml for the open string:
	db, err := sql.Open("postgres", cfg.Database.Open)
	log.Debugf("Created db connection: %v", db)
	if err != nil {
		log.Fatal(err)
	}

	var s storage.StorageBackend
	s, err = storage.NewStorageBackend(cfg)

	dbClient := dbclient.NewDBClient(db)
	i := importer.NewImporter(dbClient, s)

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Index)

	r.Handle("/api/v1/media", api.MediaHandler{Cfg: cfg, Importer: i}).Methods("POST")

	// If nothing else has matched attempt to serve a static file:
	// TODO: Point to a proper location for the data files? Use env var or config.
	r.PathPrefix("/").Handler(http.FileServer(
		http.Dir("/home/dev/go/src/github.com/dgoodwin/gophoto/public/")))

	log.Infof("Listening on port %d", cfg.APIPort)
	return &http.Server{Addr: fmt.Sprintf(":%d", cfg.APIPort), Handler: r}
}
