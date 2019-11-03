package importer

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	//log "github.com/Sirupsen/logrus"
	"github.com/dgoodwin/gophoto/pkg/api/v1/dbclient"
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
	dbClient *dbclient.DBClient
	storage  storage.StorageBackend
}

func NewImporter(dbClient *dbclient.DBClient, storage storage.StorageBackend) *Importer {
	return &Importer{dbClient: dbClient, storage: storage}
}

// ImportFilePath imports a file from the local filesystem.
func (i *Importer) ImportFilePath(filepath, checksum string) error {
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

	err = i.storage.ImportFilePath(filepath)
	if err != nil {
		return err
	}

	err = i.dbClient.CreateMedia(tm, filepath, width, height, checksum, fi.Size())
	if err != nil {
		return err
	}
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
