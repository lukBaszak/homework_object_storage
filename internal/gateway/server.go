package gateway

import (
	"context"
	"errors"
	"io"
	"log"
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
	}
}

func (o *ObjectGatewayServer) getObject(ctx context.Context, w http.ResponseWriter, id string) {
	reader, err := o.store.Get(ctx, id)

	if err != nil {
		if errors.Is(err, storage.ErrObjectNotFound) {
			log.Printf("error: %v", err)
			http.Error(w, "object not found", http.StatusNotFound)
			return
		}
		log.Printf("error2: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "failed to stream object", http.StatusInternalServerError)
		return
	}
}
