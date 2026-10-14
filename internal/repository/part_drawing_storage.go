package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// PartDrawingStorage keeps uploaded part drawing PDFs as files in one directory
// (UPLOAD_DIR/part-drawings). Each file gets a random name; part_processes.drawing_file
// holds that name. Callers pass only names checked with dto.ValidPartDrawingFile, so a
// name never leaves the directory.
type PartDrawingStorage struct {
	directory string
}

// NewPartDrawingStorage creates a PartDrawingStorage that keeps its files in directory.
func NewPartDrawingStorage(directory string) *PartDrawingStorage {
	return &PartDrawingStorage{directory: directory}
}

// Save writes content to a new file and returns its name (32 hex characters + .pdf).
func (storage *PartDrawingStorage) Save(content []byte) (string, error) {
	if err := os.MkdirAll(storage.directory, 0o750); err != nil {
		return "", fmt.Errorf("save part drawing: %w", err)
	}
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("save part drawing: %w", err)
	}
	name := hex.EncodeToString(randomBytes) + ".pdf"

	file, err := os.OpenFile(filepath.Join(storage.directory, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err != nil {
		return "", fmt.Errorf("save part drawing: %w", err)
	}
	if _, err := file.Write(content); err != nil {
		_ = file.Close()
		return "", fmt.Errorf("save part drawing: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("save part drawing: %w", err)
	}
	return name, nil
}

// Exists reports whether the drawing name is stored.
func (storage *PartDrawingStorage) Exists(name string) (bool, error) {
	_, err := os.Stat(filepath.Join(storage.directory, name))
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	default:
		return false, fmt.Errorf("find part drawing: %w", err)
	}
}

// Open opens the drawing name for reading, or returns ErrNotFound. The caller closes it.
func (storage *PartDrawingStorage) Open(name string) (*os.File, error) {
	file, err := os.Open(filepath.Join(storage.directory, name))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("open part drawing: %w", err)
	}
	return file, nil
}
