package api

import (
	"crypto/md5"
	"encoding/hex"
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
	MD5         string
}

type MediaHandler struct {
	Cfg      config.GophotoConfig
	Importer importer.Importer
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

	// Validate the MD5 checksum matches what the caller sent. For now, we won't
	// consider it an error if no MD5 checksum was provided.
	if v.MD5 != "" {
		hasher := md5.New()
		hasher.Write(v.Content)
		reqChecksum := hex.EncodeToString(hasher.Sum(nil))
		log.Debugf("Incoming data checksum: %s", reqChecksum)
		log.Debugf("Caller specified checksum: %s", v.MD5)
		if reqChecksum != v.MD5 {
			http.Error(w, "checksum mismatch", http.StatusBadRequest)
			return
		}
		log.Infoln("Checksum validated")
	} else {
		log.Warnln("No MD5 checksum specified in upload request")
	}

	h.Importer.ImportFilePath(uploadedFilePath)

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
