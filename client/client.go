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
	log.Infoln("Running client...")
	filePath := "/home/dev/Photos/2017/03/11/20170311_192038.jpg"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	v := api.Media{Name: "20170311_192038.jpg", Description: "my picture", Content: content, MD5: "789c6b4218f65047c184574dc584e19f"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8200/api/v1/media", &buf)
	if err != nil {
		// handle error
		log.Fatalln(err)
	}
	resp, err := http.DefaultClient.Do(req)
	log.Infoln(resp)
}
