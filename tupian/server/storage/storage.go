package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const UploadDir = "uploads"

func init() {
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(UploadDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("Failed to create upload directory: %v", err))
		}
	}
}

type SavedFile struct {
	UUID         string
	Filename     string
	OriginalName string
	FilePath     string
	FileSize     int64
	ContentType  string
}

func SaveFile(file *multipart.FileHeader) (*SavedFile, error) {
	if file == nil {
		return nil, errors.New("file is nil")
	}

	fileUUID := uuid.New().String()
	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFilename := fileUUID + ext
	filePath := filepath.Join(UploadDir, newFilename)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	contentType := file.Header.Get("Content-Type")

	return &SavedFile{
		UUID:         fileUUID,
		Filename:     newFilename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     fileSize,
		ContentType:  contentType,
	}, nil
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func GetFilePath(filename string) string {
	return filepath.Join(UploadDir, filename)
}
