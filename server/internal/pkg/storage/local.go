package storage

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage stores files on the local filesystem.
type LocalStorage struct {
	BasePath string
}

// Save writes data to BasePath/{fileID}_{filename} and returns the stored filename.
func (s *LocalStorage) Save(fileID, filename string, data io.Reader) (string, error) {
	if err := os.MkdirAll(s.BasePath, 0755); err != nil {
		return "", fmt.Errorf("create storage dir: %w", err)
	}
	storedName := fileID + "_" + filename
	dest := filepath.Join(s.BasePath, storedName)
	f, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, data); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return storedName, nil
}

// Get opens the file identified by fileID (prefix match) and returns a reader and MIME type.
func (s *LocalStorage) Get(fileID string) (io.ReadCloser, string, error) {
	entries, err := os.ReadDir(s.BasePath)
	if err != nil {
		return nil, "", fmt.Errorf("read storage dir: %w", err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), fileID+"_") {
			path := filepath.Join(s.BasePath, e.Name())
			f, err := os.Open(path)
			if err != nil {
				return nil, "", err
			}
			ext := filepath.Ext(e.Name())
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
			return f, mimeType, nil
		}
	}
	return nil, "", fmt.Errorf("file not found: %s", fileID)
}

// Delete removes the file identified by fileID.
func (s *LocalStorage) Delete(fileID string) error {
	entries, err := os.ReadDir(s.BasePath)
	if err != nil {
		return fmt.Errorf("read storage dir: %w", err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), fileID+"_") {
			return os.Remove(filepath.Join(s.BasePath, e.Name()))
		}
	}
	return nil // not found is not an error
}
