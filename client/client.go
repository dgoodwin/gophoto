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

	// TODO: calculate sha1, fix filename
	v := api.Media{Name: "20170311_192038.jpg", Description: "my picture", Content: content, Checksum: "469c00f6204d08f362e2bb16c56396de1db4a14f"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		log.WithError(err).Fatal("error encoding file")
	}
	// TODO: configurable endpoint
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
