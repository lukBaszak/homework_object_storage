package gateway

import (
	"context"
	"fmt"
	"main/internal/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGETFile(t *testing.T) {

	defaultDir := "testing"
	filename := "default"
	initialData := []byte("Hello world")
	store := storage.FileSystemStore{}

	server := NewObjectGatewayServer(&store)

	if err := store.Setup(context.Background(), defaultDir); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	t.Run("Retrieve existing object with correct path", func(t *testing.T) {
		path := filepath.Join(defaultDir, filename)
		if err := os.WriteFile(path, initialData, 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/object/%s", filename), nil)
		response := httptest.NewRecorder()

		server.router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "Hello world")
	})

	t.Run("Retrieve non-existing object with correct path", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/object/%s", "incorrect"), nil)
		response := httptest.NewRecorder()

		server.router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("Retrieve with incorrect path", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/object/%s", strings.Repeat("0", 33)), nil)
		response := httptest.NewRecorder()

		server.router.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
