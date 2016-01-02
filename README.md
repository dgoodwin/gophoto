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

Build a quick and dirty webapp which can import files from disk into a database, generate some thumbnails, expose a REST API for those photos and their thumbnails, and then see how hard it would be to build a responsive utilitarian Javascript UI to view them.

# Running In Docker

Preferred method of deployment is within a Docker container. A Docker compose configuration is provided to quickly build the container from source and launch with an accompanying PostgreSQL database. docker-compose can be used to quickly teardown and recreate this environment.

TODO: This is temporary but make sure you edit docker-compose.yml to point to a directory on your system with some photos to import.

```
$ sudo docker-compose build
$ sudo docker-compose up -d
$ sudo docker-compose run --rm gophoto goose up
```

Use the `docker logs` command on each container to view activity and logging.

Use `docker-compose stop && docker-compose rm` to destroy the environment completely.

Use `docker exec -ti gophoto_db_1 psql -U postgres gophoto` to access the database manually.

# Running From Source

You probably don't want to do this...

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



