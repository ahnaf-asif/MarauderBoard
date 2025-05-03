package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open uploaded file: %w", err)
	}
	defer src.Close()

	destDir := "./public/uploads"

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create upload directory: %w", err)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	destPath := filepath.Join(destDir, filename)

	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	publicURI := "/static/uploads/" + filename
	return publicURI, nil
}
