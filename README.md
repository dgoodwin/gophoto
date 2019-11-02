# About

If this project were to not collapse in a sea of apathy and disinterest, which it will, it would eventually become a self-hosted open source alternative to Google Photos.

  * Built to handle a lifetime's worth of photos.
  * Able to browse quickly and easily by date or albums.
  * Content stored in original quality.
  * Strong sharing options including private links, and a public feed for things you flag as such.
  * Rich web UI for organizing.
  * Possibly video support someday.
  * Able to upload (or auto-upload) from Android and Shotwell. (or with a CLI binary)
  * Federation via defined open standards.

# Current Goals

Quick and dirty webapp exposing a REST API which can upload and query photos.
Thumbnails generated on upload, and everything should be stored in S3.
(possibly bring your own credentials)

Minimal client binary that can search a directory recursively and upload all
the images it finds, if they haven't been uploaded already.

Minimal JS UI that renders image thumbnails, and loads more as you scroll.
Hopefully re-use a react photo gallery library.

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

You probably don't want to do this unless you're looking to actually hack on
the code and find the docker method above too slow. (which it probably is for
rapid development)

Ensure you have your go workspace and GOPATH environment variable setup correctly.

Compile the project:

```
$ mkdir -p $GOPATH/src/github.com/dgoodwin/gophoto/
$ git clone git@github.com:dgoodwin/gophoto.git $GOPATH/src/github.com/dgoodwin/gophoto/
$ cd $GOPATH/src/github.com/dgoodwin/gophoto/
$ go get
$ go install ./...
```

Install goose for managing the database schema, launch postgresql in a
container (or use another server, but you will need to update dbconf.yml
appropriately), create the database and populate schema:

```
$ go get bitbucket.org/liamstask/goose/cmd/goose
$ docker run --name postgres-dev -e POSTGRES_PASSWORD=random --restart=always -p 15432:5432 -d postgres
$ PGPASSWORD="random" createdb -U postgres -h localhost -p 15432 gophoto
$ goose -env local up
```

Edit a copy of the gophoto-docker.yml for configuration:

```
$ cat gophoto-local.yml
storage:
    backend: fileSystem
    path: ./storage/
assetspath: ./public
importpath: /home/dev/Photos/2015/12/
database:
    open: user=postgres dbname=gophoto sslmode=disable host=localhost port=15432 password=random
```

And finally launch the app:

```
$ gophoto serve gophoto-local.yml
```

# Running Tests

```
$ go test ./...
```



