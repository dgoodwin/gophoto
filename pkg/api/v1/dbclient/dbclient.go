package dbclient

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
)

func NewDBClient(db *sql.DB) *DBClient {
	return &DBClient{db: db}
}

type DBClient struct {
	db *sql.DB
}

// TODO: refactor to take a model object?
func (dbc *DBClient) CreateMedia(taken time.Time, filename string, res_x int, res_y int, checksum string, size int64) error {
	var newPhotoId int
	stmt, err := dbc.db.Prepare("INSERT INTO media(created, uploaded, filename, url, checksum, res_x, res_y, size) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(taken, time.Now(), filename, filename, checksum, res_x, res_y, size).Scan(&newPhotoId)
	if err != nil {
		return err
	}
	log.Infof("Created new photo in db: %d", newPhotoId)
	return nil
}
