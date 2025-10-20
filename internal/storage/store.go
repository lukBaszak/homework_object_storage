package storage

import (
	"context"
	"io"
)

type Store interface {
	Setup(ctx context.Context, defaultLocation string) error
	Get(ctx context.Context, fileName string) (io.ReadCloser, error)
}
