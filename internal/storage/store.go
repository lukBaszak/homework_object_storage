package storage

import (
	"context"
	"io"
)

type Store interface {
	Setup(ctx context.Context, defaultLocation string) error
	Get(ctx context.Context, fileName string) (io.ReadCloser, error)
	Put(ctx context.Context, fileName string, reader io.Reader, size int64) error
}
