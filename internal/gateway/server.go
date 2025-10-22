package gateway

import (
	"context"
	"errors"
	"io"
	"main/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

type ObjectGatewayServer struct {
	store  storage.Store
	router *mux.Router
}

func NewObjectGatewayServer(store storage.Store) *ObjectGatewayServer {
	s := new(ObjectGatewayServer)

	s.store = store

	router := mux.NewRouter()
	router.HandleFunc("/object/{id:[A-Za-z0-9]{1,32}}", s.ObjectHandler)

	s.router = router

	return s
}

func (o *ObjectGatewayServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.router.ServeHTTP(w, r)
}

func (o *ObjectGatewayServer) ObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		o.getObject(ctx, w, vars["id"])
	case http.MethodPut:
		o.putObject(ctx, w, r.Body, vars["id"], r.ContentLength)
	}
}

func (o *ObjectGatewayServer) getObject(ctx context.Context, w http.ResponseWriter, id string) {
	reader, err := o.store.Get(ctx, id)

	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			http.Error(w, "object not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "failed to stream object", http.StatusInternalServerError)
		return
	}
}

func (o *ObjectGatewayServer) putObject(ctx context.Context, w http.ResponseWriter, readCloser io.ReadCloser, id string, size int64) {
	status, err := o.store.Put(ctx, id, readCloser, size)

	if err != nil {
		http.Error(w, "failed to put object", http.StatusInternalServerError)
	}

	switch status {
	case storage.PutCreated:
		w.WriteHeader(http.StatusCreated)
	case storage.PutOverwritten:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "failed to stream object", http.StatusInternalServerError)
	}
}
