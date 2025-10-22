package main

import (
	"context"
	"log"
	"main/internal/docker"
	"main/internal/gateway"
	"main/internal/storage"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/client"
)

const (
	minioDockerLabelKey   = "service-type"
	minioDockerLabelValue = "minio-worker-node"
	dockerNetworkName     = "minio-net"
	maxRetries            = 3
)

func main() {

	bucketName := os.Getenv("BUCKET_NAME")

	if bucketName == "" {
		bucketName = "default"
	}

	ctx := context.Background()

	log.Println("initializing docker client...")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatal("couldn't initialize docker client...")
	}

	minioConfig, err := docker.GetMinioConfig(ctx, cli, minioDockerLabelKey, minioDockerLabelValue, dockerNetworkName)

	if err != nil {
		log.Fatalf("couldn't retrieve minio config: %v", err)
	}

	log.Println("minio config retrieved successfully")

	nodes := storage.NewNodesConfig(minioConfig)
	store := storage.MinioStore{Nodes: nodes}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = store.Setup(context.Background(), bucketName)
		if err == nil {
			log.Printf("successfully set up bucket - %q on attempt %d", bucketName, attempt)
			break
		}

		log.Printf("attempt %d/%d failed: %v", attempt, maxRetries, err)

		if attempt == maxRetries {
			log.Fatalf("couldn't setup initial buckets")
		}

		delay := 3 * time.Second
		log.Printf("retrying in %v...", delay)
		time.Sleep(delay)
	}

	server := gateway.NewObjectGatewayServer(&store)

	log.Fatal(http.ListenAndServe(":3000", server))
}
