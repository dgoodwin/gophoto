package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type Media struct {
	Name        string
	Description string
	Content     []byte
}

func MediaPost(w http.ResponseWriter, r *http.Request) {
	reqLog := log.WithFields(log.Fields{
		"method": r.Method,
		"len":    r.ContentLength,
	})

	reqLog.Infoln("MediaPost")
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
	reqLog.Infof("Uploaded file: %s", v.Name)

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
