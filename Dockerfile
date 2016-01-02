# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:onbuild

# Install goose for managing the database.
RUN go get bitbucket.org/liamstask/goose/cmd/goose

EXPOSE 8080
