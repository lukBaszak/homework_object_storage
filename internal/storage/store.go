package storage

type Store interface {
	Get(dir, file string) ([]byte, error)
}
