package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	FileNotFoundException = "couldn't find given file"
)

type FileSystemStore struct {
	defaultLocation string
}

func (f *FileSystemStore) Setup(ctx context.Context, defaultLocation string) error {
	err := os.Mkdir(defaultLocation, 0755)

	if err != nil && !os.IsExist(err) {
		return err
	}

	f.defaultLocation = defaultLocation

	return nil
}

func (f *FileSystemStore) Get(ctx context.Context, fileName string) ([]byte, error) {
	path := filepath.Join(f.defaultLocation, fileName)

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("%s: %w", FileNotFoundException, err)
	}

	file, _ := os.ReadFile(path)

	return file, nil
}
