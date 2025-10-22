package main

import (
	"context"
	"log"
	"main/internal/docker"
	"main/internal/gateway"
	"main/internal/storage"
	"net/http"

	"github.com/docker/docker/client"
)

func main() {

	ctx := context.Background()
	log.Println("initializing docker client...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("couldn't initialize docker client...")
	}
	minioConfig, err := docker.GetMinioConfig(ctx, cli, "service-type", "minio-worker-node", "minio-net")

	if err != nil {
		log.Fatalf("couldn't retrieve minio config: %v", err)
	}

	log.Printf("minio config retrieved successfully: %s", minioConfig)

	nodes := storage.NewNodesConfig(minioConfig)
	store := storage.MinioStore{Nodes: nodes}
	store.Setup(ctx, "default")

	server := gateway.NewObjectGatewayServer(&store)

	log.Fatal(http.ListenAndServe(":3000", server))
}
