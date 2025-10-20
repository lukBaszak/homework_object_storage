package storage

import "context"

type Store interface {
	Setup(ctx context.Context, defaultLocation string) error
	Get(ctx context.Context, fileName string) ([]byte, error)
}
