package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type HashRingMember string

func (h HashRingMember) String() string {
	return string(h)
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

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
		log.Println("client initialization with", cfg.Endpoint, "in progress...")

		cli, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false,
		})

		if err != nil {
			log.Printf("client initilization failed with %s", cfg.Id)
			continue
		}

		log.Printf("clients initilization succeded.")

		nodes[cfg.Id] = Node{
			client: cli,
		}
	}
	return nodes
}

type MinioStore struct {
	HashRing      *consistent.Consistent
	Nodes         map[string]Node
	DefaultBucket string
}

func (m *MinioStore) Setup(ctx context.Context, defaultLocation string) error {
	cfg := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    271,
		ReplicationFactor: 20,
		Load:              1.25,
	}
	m.HashRing = consistent.New(nil, cfg)

	for id, node := range m.Nodes {
		err := CreateBucket(ctx, defaultLocation, node, id)

		if err != nil {
			return err
		}
		m.HashRing.Add(HashRingMember(id))
	}

	m.DefaultBucket = defaultLocation

	return nil
}

func (m *MinioStore) Get(ctx context.Context, filenameId string) (io.ReadCloser, error) {
	key := []byte(filenameId)
	owner := m.HashRing.LocateKey(key)

	obj, err := m.Nodes[owner.String()].client.GetObject(ctx, m.DefaultBucket, filenameId, minio.GetObjectOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to get object from node %s: %w", owner, err)
	}

	_, err = obj.Stat()

	if err != nil {
		var resp minio.ErrorResponse
		if errors.As(err, &resp) && resp.Code == "NoSuchKey" {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("read object: %w", err)
	}
	return obj, nil
}

func (m *MinioStore) Put(ctx context.Context, fileNameId string, reader io.Reader, size int64) error {
	key := []byte(fileNameId)
	owner := m.HashRing.LocateKey(key)

	info, err := m.Nodes[owner.String()].client.PutObject(ctx, m.DefaultBucket, fileNameId, reader, size, minio.PutObjectOptions{})

	if err != nil {
		var resp minio.ErrorResponse
		if errors.As(err, &resp) {
			switch resp.Code {
			case "NoSuchBucket":
				return fmt.Errorf("bucket %s does not exist on node %s: %w", m.DefaultBucket, owner, err)
			default:
				return fmt.Errorf("minio error [%s] from node %s: %w", resp.Code, owner, err)
			}
		}
		return fmt.Errorf("failed to put object on node %s: %w", owner, err)
	}

	log.Printf("Object %s uploaded to node %s, size=%d, etag=%s", fileNameId, owner, info.Size, info.ETag)

	return nil
}

func CreateBucket(ctx context.Context, bucketName string, node Node, id string) error {
	exists, err := node.client.BucketExists(ctx, bucketName)

	if err != nil {
		log.Printf("there was an error during initial setup of %s node: %v", id, err)
	}

	if !exists {
		if err := node.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("there was an error during initial bucket creation in %s node: %v", id, err)
		}
	}
	return nil
}
