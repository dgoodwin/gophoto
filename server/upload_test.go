package server

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	log "github.com/Sirupsen/logrus"
)

// Will store a standard test suite of photos for upload tests here.
// This code will sync them to a local directory if they don't already
// exist.
var testSuitePhotosLocation = "http://files.rm-rf.ca/gophoto-tests/"

var testSuitePhotos = map[string]struct {
	filename string
	md5sum   string
	date     string
}{
	"911": {
		filename: "20160820_152200.jpg",
		md5sum:   "376fc377b5ec5e60bc5b7a29ccd53fd2",
		date:     "2016:08:20 15:22:00",
	},
	"911-no-exif": {
		filename: "20160820_152200_no_exif.jpg",
		md5sum:   "708f8922ce3585328b600a4b6dce1b0f",
		date:     "",
	},
}

func syncTestPhotos() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	testSuiteDir := filepath.Join(cwd, "photo-test-suite")
	log.Infof("Photo test suite: %s", testSuiteDir)

	if _, err := os.Stat(testSuiteDir); os.IsNotExist(err) {
		log.Infoln("Creating photo test suite dir:", testSuiteDir)
		err := os.Mkdir(testSuiteDir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, p := range testSuitePhotos {
		expectedPath := filepath.Join(testSuiteDir, p.filename)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			u, _ := url.Parse(testSuitePhotosLocation)
			u.Path = path.Join(u.Path, p.filename)

			log.Infoln("Downloading:", u.String())
			resp, err := http.Get(u.String())
			if err != nil {
				log.Fatal(err)
			} else {
				defer resp.Body.Close()
				out, err := os.Create(expectedPath)
				if err != nil {
					log.Fatal(err)
				}
				defer out.Close()
				io.Copy(out, resp.Body)
			}
		} else {
			log.Infoln("Test photo already exists:", expectedPath)
		}

		// Make sure the checksum matches:
		hasher := md5.New()
		f, err := os.Open(expectedPath)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(hasher, f); err != nil {
			log.Fatal(err)
		}

		actualChecksum := hex.EncodeToString(hasher.Sum(nil))
		if actualChecksum != p.md5sum {
			log.Fatalf("checksum mismatch on %s: expected=%s got=%s", expectedPath, p.md5sum, actualChecksum)
		} else {
			log.Infoln("Checksum validated:", p.md5sum)
		}
	}
}

func TestSomething(t *testing.T) {
	syncTestPhotos()
	// Bring up a server
	/*
		for _, p := range testSuitePhotos {
		}
	*/
}
