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

func (f *FileSystemStore) Put(ctx context.Context, fileName string, reader io.Reader, size int64) (int, error) {
	filePath := filepath.Join(f.defaultLocation, fileName)

	fileExisted := fileExists(filePath)
	file, err := os.Create(filePath)

	if err != nil {
		return PutError, fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, reader)

	if err != nil {
		return PutError, fmt.Errorf("copy data: %w", err)
	}

	if written != size {
		return PutError, io.ErrShortWrite
	}

	if fileExisted {
		return PutOverwritten, nil
	}

	return PutCreated, nil
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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
