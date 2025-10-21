package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func (f *FileSystemStore) Put(ctx context.Context, fileName string, reader io.Reader, size int64) error {
	filePath := filepath.Join(f.defaultLocation, fileName)

	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	written, err := io.Copy(file, reader)

	if err != nil {
		return err
	}

	if written != size {
		return io.ErrShortWrite
	}

	return nil
}

func (f *FileSystemStore) Get(ctx context.Context, fileName string) (io.ReadCloser, error) {
	path := filepath.Join(f.defaultLocation, fileName)

	file, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
