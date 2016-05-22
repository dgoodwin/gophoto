package importer

import (
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

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

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

func saveMetadata(db *sql.DB, filename string, res_x int, res_y int, size int64) error {
	var newPhotoId int
	stmt, err := db.Prepare("INSERT INTO media(filename, url, res_x, res_y, size) VALUES($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return err
	}
	fmt.Printf("stmt prepared\n")
	err = stmt.QueryRow(filename, filename, res_x, res_y, size).Scan(&newPhotoId)
	if err != nil {
		return err
	}
	fmt.Printf("Created new photo: %d\n", newPhotoId)
	return nil
}
