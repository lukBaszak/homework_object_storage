package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("Retrieve file in existing path", func(t *testing.T) {

		tmpDir := t.TempDir()
		filename := "testFile"
		initialData := []byte("Hello world")

		createTempFile(t, tmpDir, filename, initialData)

		store := FileSystemStore{}
		got, _ := store.Get(tmpDir, filename)

		assertData(t, initialData, got)

	})

	t.Run("Retrieve file in non-existing path", func(t *testing.T) {
		tmpDir := t.TempDir()
		filename := "testFile"
		initialData := []byte("Hello world")

		createTempFile(t, tmpDir, filename, initialData)

		store := FileSystemStore{}
		_, err := store.Get(tmpDir, "non-existing")

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
