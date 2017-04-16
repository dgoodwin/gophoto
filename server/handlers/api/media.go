package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dgoodwin/gophoto/config"
	"github.com/dgoodwin/gophoto/server/importer"

	log "github.com/Sirupsen/logrus"
)

type Media struct {
	Name        string
	Description string
	Content     []byte
}

type MediaHandler struct {
	Cfg      config.GophotoConfig
	Importer importer.Importer
}

// TODO: Probably a lot of work to come here. Currently this accepts JSON with
// b64 encoded content. Need to consider uploads from browser or curl, chunking,
// etc.
func (h MediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"len":    r.ContentLength,
	}).Infoln("handling request")

	if r.Body == nil {
		http.Error(w, "no body in request", http.StatusBadRequest)
		return
	}

	/*
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	var v Media
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("Uploaded file: %s", v.Name)

	// Create a temporary directory for this upload in the working dir. Used to
	// store the incoming file and generate thumbnails before we push to permanent
	// storage.
	tempUploadDir, err := ioutil.TempDir(h.Cfg.WorkDirPath, "upload-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("Created temporary upload dir: %s", tempUploadDir)
	defer func() {
		log.Infof("Cleaning up temporary upload dir: %s", tempUploadDir)
		os.RemoveAll(tempUploadDir)
	}()

	uploadedFilePath := filepath.Join(tempUploadDir, v.Name)
	log.Infof("Uploaded file: %s", uploadedFilePath)
	if err := ioutil.WriteFile(uploadedFilePath, v.Content, 0644); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		path := mux.Params(r).ByName("path")
		err = c.bucket.Put(path, content, r.Header.Get("Content-Type"), s3.PublicRead)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		url := c.bucket.URL(path)
		w.Header().Set("Location", url)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"url":"%s"}`, url)
	*/
	return
}
