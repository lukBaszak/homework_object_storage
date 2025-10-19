package storage

type Store interface {
	Setup(defaultLocation string) error
	Get(file string) ([]byte, error)
}
