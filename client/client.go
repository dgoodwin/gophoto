package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dgoodwin/gophoto/server/handlers/api"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RunSync(cmd *cobra.Command, args []string) {
	log.Infoln("Running client hi...")
	filePath := "/home/dev/Photos/2017/03/11/20170311_192038.jpg"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.WithError(err).Fatal("error reading file")
	}

	// TODO: calculate md5, fix filename
	v := api.Media{Name: "20170311_192038.jpg", Description: "my picture", Content: content, MD5: "789c6b4218f65047c184574dc584e19f"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		log.WithError(err).Fatal("error encoding file")
	}
	req, err := http.NewRequest("POST", "http://localhost:8201/api/v1/media", &buf)
	if err != nil {
		// handle error
		log.WithError(err).Fatal("error making API request")
	}
	log.Infof("request: %v", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Fatal("request error")
	}
	log.Infof("response: %v", resp)
}
