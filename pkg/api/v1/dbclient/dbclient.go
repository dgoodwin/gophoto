package dbclient

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"

	api "github.com/dgoodwin/gophoto/pkg/api/v1"
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

func (dbc *DBClient) ListMedia() ([]*api.Media, error) {
	mediaResults := []*api.Media{}
	rows, err := dbc.db.Query("select created, uploaded, filename, id, url, checksum, res_x, res_y, size from media")
	if err != nil {
		return mediaResults, err
	}
	defer rows.Close()
	var (
		created  time.Time
		uploaded time.Time
		filename string
		id       int
		url      string
		checksum string
		resX     int
		resY     int
		size     int
	)
	for rows.Next() {
		err := rows.Scan(&created, &uploaded, &filename, &id, &url, &checksum, &resX, &resY, &size)
		if err != nil {
			return mediaResults, err
		}
		mediaResults = append(mediaResults,
			&api.Media{
				Created:  created,
				Uploaded: uploaded,
				FileName: filename,
				ID:       id,
				URL:      url,
				Checksum: checksum,
				Resolution: api.Resolution{
					X: resX,
					Y: resY,
				},
				Size: size,
			})
	}
	err = rows.Err()
	if err != nil {
		return mediaResults, err
	}
	return mediaResults, nil
}
