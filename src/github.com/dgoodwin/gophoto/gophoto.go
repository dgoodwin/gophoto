package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && isImage(path) {
		fmt.Printf("Visited: %s\n", path)
		//fmt.Printf("  %d %s %s\n", f.Size(), f.Mode(), f.IsDir())
		width, height := getImageDimension(path)
		fmt.Println("  Width:", width, "  Height:", height)
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

func main() {
	fmt.Printf("Hello, world.\n")
	importDir := "/home/dev/Photos/2015/11"
	err := filepath.Walk(importDir, visit)
	fmt.Printf("Walk returned: %v\n", err)
}
