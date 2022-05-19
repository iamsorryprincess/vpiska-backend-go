package storage

type FileStorage interface {
	Get(id string) ([]byte, error)
	Upload(id string, data []byte) error
	Delete(id string) error
}
