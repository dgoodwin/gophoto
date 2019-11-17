package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dgoodwin/gophoto/config"
	api "github.com/dgoodwin/gophoto/pkg/api/v1"
	"github.com/dgoodwin/gophoto/pkg/api/v1/dbclient"
	"github.com/dgoodwin/gophoto/server/importer"

	log "github.com/Sirupsen/logrus"
)

type MediaHandler struct {
	Cfg      config.GophotoConfig
	Importer *importer.Importer
	DBClient *dbclient.DBClient
}

/*
Planned upload API:

- POST /media
  - accepts a list of media objects minus actual data
  - requests upload permission
  - returns a list of media entities including href where the data can be PUT
  - validates checksum is unique, won't let you re-upload (in future may have override)
  - generates a short unique ID for URL
  - stores record in database
  - unused records in db will be cleaned at some point in future if overdue for upload
- PUT /media/ABCDEFGH
  - now includes data
  - can handle chunked uploads if user desires (future feature)
  - can be either json with b64 encoded fields, or raw HTTP PUT with headers for metadata
*/

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

	switch method := r.Method; method {
	case http.MethodPost:
		h.serveCreate(w, r)
		return
	case http.MethodGet:
		h.serveGet(w, r)
		return
	default:
		msg := fmt.Sprintf("unsupported http method: %s", method)
		log.Warn(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func (h MediaHandler) serveCreate(w http.ResponseWriter, r *http.Request) {

	// TODO: error if any unexpected fields are set or unset.

	/*
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	var v api.Media
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof("Uploaded file: %s", v.FileName)

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

	uploadedFilePath := filepath.Join(tempUploadDir, v.FileName)
	log.Infof("Uploaded file: %s", uploadedFilePath)
	if err := ioutil.WriteFile(uploadedFilePath, v.Content, 0644); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the checksum matches what the caller sent. For now, we won't
	// consider it an error if no checksum was provided.
	if v.Checksum != "" {
		h := sha1.New()
		h.Write(v.Content)
		reqChecksum := hex.EncodeToString(h.Sum(nil))
		log.Debugf("Incoming data checksum: %s", reqChecksum)
		log.Debugf("Caller specified checksum: %s", v.Checksum)
		if reqChecksum != v.Checksum {
			http.Error(w, "checksum mismatch", http.StatusBadRequest)
			return
		}
		log.Infoln("Checksum validated")
	} else {
		log.Warnln("No checksum specified in upload request")
	}

	if err := h.Importer.ImportFilePath(uploadedFilePath, v.Checksum); err != nil {
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

func (h MediaHandler) serveGet(w http.ResponseWriter, r *http.Request) {

	media, err := h.DBClient.ListMedia()
	if err != nil {
		log.WithError(err).Error("error listing media")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonOutput, err := json.Marshal(media)
	if err != nil {
		log.WithError(err).Error("error encoding media array as json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonOutput)

	return
}
