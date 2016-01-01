# About

If this project were to not collapse in a sea of apathy and disinterest, which it will, it would eventually become a self-hosted open source alternative to Google Photos.

  * Built to handle a lifetime's worth of photos.
  * Able to browse quickly and easily by date or albums.
  * Content stored in original quality.
  * Strong sharing options including private links.
  * Rich web UI for organizing.
  * Possibly eventual video support.
  * Able to upload (or auto-upload) from Android and Shotwell.

# Current Status

I have no idea what I'm doing.

# Running In Docker

Gophoto is intended to be run within a Docker container. Docker compose can be used to quickly spin up the application and a postgresql database, and teardown/recreate as necessary.

```
docker build -t gophoto .
```

# Running From Source

Ensure you have GOPATH environment variable setup correctly.

```
$ mkdir -p $GOPATH/src/github.com/dgoodwin/gophoto/.
$ git clone git@github.com:dgoodwin/gophoto.git $GOPATH/src/github.com/dgoodwin/gophoto/
$ cd $GOPATH/src/github.com/dgoodwin/gophoto/
$ go get
$ go install github.com/dgoodwin/gophoto
$ gophoto [CONFIGPATH]
```

# Running Tests

```
$ go test ./...
```



