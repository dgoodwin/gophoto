package importer

import (
	"database/sql"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgoodwin/gophoto/server/storage"
	"github.com/rwcarlsen/goexif/exif"
)

/*
func CheckFile(path string, f os.FileInfo, db *sql.DB) error {
	if !f.IsDir() && isImage(path) {
		fmt.Printf("Visited: %s\n", path)
		//fmt.Printf("  %d %s %s\n", f.Size(), f.Mode(), f.IsDir())
		width, height := getImageDimension(path)
		fmt.Println("  Width:", width, "  Height:", height)
		err := saveMetadata(db, path, width, height, f.Size())
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Skipping non-image: %s\n", path)
	}
	return nil
}

func isImage(path string) bool {
	imageExtensions := map[string]bool{
		".jpg": true,
		".JPG": true,
	}
	if imageExtensions[filepath.Ext(path)] {
		return true
	}
	return false
}
*/

func getImageDimension(filepath string) (int, int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	image, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

// Importer handles an incoming file being uploaded, orchestrates thumbnail
// generation, stores metadata in the database, and forwards on to it's final
// storage backend.
type Importer struct {
	DB      *sql.DB
	Storage storage.StorageBackend
}

// ImportFilePath imports a file from the local filesystem.
func (i Importer) ImportFilePath(filepath string) error {
	fi, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	width, height, err := getImageDimension(filepath)
	if err != nil {
		return err
	}

	exif, err := getExifData(filepath)
	if err != nil {
		return err
	}
	tm, _ := exif.DateTime()

	err = i.Storage.ImportFilePath(filepath)
	if err != nil {
		return err
	}

	err = i.saveMetadata(tm, filepath, width, height, fi.Size())
	if err != nil {
		return err
	}
	return nil
}

func (i Importer) saveMetadata(taken time.Time, filename string, res_x int, res_y int, size int64) error {
	var newPhotoId int
	stmt, err := i.DB.Prepare("INSERT INTO media(created, uploaded, filename, url, res_x, res_y, size) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	if err != nil {
		return err
	}
	err = stmt.QueryRow(taken, time.Now(), filename, filename, res_x, res_y, size).Scan(&newPhotoId)
	if err != nil {
		return err
	}
	log.Infof("Created new photo in db: %d", newPhotoId)
	return nil
}

func getExifData(filepath string) (*exif.Exif, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// May want to just go straight to the CLI if this isn't reliable enough.
	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	return x, nil
}
