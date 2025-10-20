package storage

import (
	"context"
	"fmt"
	"io"
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

func (f *FileSystemStore) Get(ctx context.Context, fileName string) (io.ReadCloser, error) {
	path := filepath.Join(f.defaultLocation, fileName)

	file, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %w", FileNotFoundException, err)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
