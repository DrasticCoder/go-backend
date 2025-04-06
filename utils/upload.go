package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const UploadDir = "uploads" // ensure this directory exists

func SaveMedia(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Generate a unique filename
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	path := filepath.Join(UploadDir, filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return path, nil
}
