package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("Setup with default directory", func(t *testing.T) {
		defaultDir := "testing"

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

		store := FileSystemStore{defaultLocation: "testing"}

		filename := "testFile"
		initialData := []byte("Hello world")

		createTempFile(t, store.defaultLocation, filename, initialData)

		got, _ := store.Get(context.Background(), filename)

		assertData(t, initialData, got)

	})

	t.Run("Retrieve file in non-existing path", func(t *testing.T) {

		store := FileSystemStore{defaultLocation: "testing"}

		filename := "testFile"
		initialData := []byte("Hello world")

		createTempFile(t, store.defaultLocation, filename, initialData)

		_, err := store.Get(context.Background(), "non-existing")

		assertError(t, err)
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
