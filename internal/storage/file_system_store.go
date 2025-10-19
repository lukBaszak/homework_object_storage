package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	FileNotFoundException = "couldn't find given file"
)

type FileSystemStore struct {
}

func (f FileSystemStore) Get(dir, fileName string) ([]byte, error) {
	path := filepath.Join(dir, fileName)

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("%s: %w", FileNotFoundException, err)
	}

	file, _ := os.ReadFile(path)

	return file, nil
}
