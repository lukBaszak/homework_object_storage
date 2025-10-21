package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	defaultDir := "testing"
	filename := "testFile"
	initialData := []byte("Hello world")

	t.Run("Setup with default directory", func(t *testing.T) {
		store := FileSystemStore{}
		err := store.Setup(context.Background(), defaultDir)

		if err != nil {
			t.Errorf("There was an error trying to setup directory")
		}

		if _, err = os.Stat(defaultDir); os.IsNotExist(err) {
			t.Errorf("Directory %s does not exist after setup", defaultDir)
		}
	})

	t.Run("Retrieve file in existing path", func(t *testing.T) {

		store := FileSystemStore{defaultLocation: defaultDir}

		createTempFile(t, store.defaultLocation, filename, initialData)

		reader, _ := store.Get(context.Background(), filename)
		got, err := io.ReadAll(reader)

		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		assertData(t, initialData, got)
	})

	t.Run("Retrieve file in non-existing path", func(t *testing.T) {

		store := FileSystemStore{defaultLocation: defaultDir}

		createTempFile(t, store.defaultLocation, filename, initialData)

		_, err := store.Get(context.Background(), "non-existing")

		assertError(t, err)
	})

	t.Run("Put file with chosen filename", func(t *testing.T) {

		store := FileSystemStore{defaultLocation: defaultDir}

		store.Setup(context.Background(), defaultDir)

		err := store.Put(context.Background(), filename, bytes.NewReader(initialData), int64(len(initialData)))

		if err != nil {
			t.Errorf("couldn't create file: %v", err)
		}

		reader, _ := store.Get(context.Background(), filename)
		got, _ := io.ReadAll(reader)

		assertData(t, initialData, got)
	})
}

func createTempFile(t testing.TB, dir, filename string, initialData []byte) {
	path := filepath.Join(dir, filename)

	err := os.WriteFile(path, initialData, 0644)
	if err != nil {
		t.Fatalf("Couldn't create file in %s, because of %v", path, err)
	}
}

func assertData(t testing.TB, want, got []byte) {
	t.Helper()
	if !bytes.Equal(want, got) {
		t.Errorf("Want: %v, Got: %v", string(want), string(got))
	}
}

func assertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Should throw error")
	}
}
