package storage

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Id        string
	Endpoint  string
	AccessKey string
	SecretKey string
}

type Node struct {
	client *minio.Client
}

func NewNodesConfig(minioConfigs []MinioConfig) map[string]Node {
	nodes := make(map[string]Node)

	for _, cfg := range minioConfigs {
		log.Println("Client initialization with", cfg.Endpoint, "in progress...")

		cli, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false,
		})

		if err != nil {
			log.Printf("Client initilization failed with %s", cfg.Id)
			continue
		}

		log.Printf("Clients initilization succeded.")

		nodes[cfg.Id] = Node{
			client: cli,
		}
	}
	return nodes
}

type MinioStore struct {
	Nodes map[string]Node
}

func (m *MinioStore) Setup(context context.Context, defaultLocation string) error {

	for id, node := range m.Nodes {
		err := CreateBucket(context, defaultLocation, node, id)
	}
	return nil
}

func (m *MinioStore) Get(context context.Context, file string) ([]byte, error) {
	panic("implement me")
}

func CreateBucket(context context.Context, bucketName string, node Node, id string) error {
	exists, err := node.client.BucketExists(context, bucketName)

	if err != nil {
		log.Printf("There was an error during initial setup of %s node: %v", id, err)
	}

	if !exists {
		if err := node.client.MakeBucket(context, bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("There was an error during initial bucket creation in %s node: %v", id, err)
		}
	}
	return nil
}
