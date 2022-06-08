package storage

import (
	"os"
	"sync"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
)

type localFileStorage struct {
	path    string
	mutex   sync.RWMutex
	mutexes map[string]*sync.Mutex
}

func NewLocalFileStorage(path string) (FileStorage, error) {
	err := initFileDir(path)

	if err != nil {
		return nil, err
	}

	return &localFileStorage{
		path:    path,
		mutex:   sync.RWMutex{},
		mutexes: make(map[string]*sync.Mutex),
	}, nil
}

func (s *localFileStorage) Get(id string) ([]byte, error) {
	file, err := os.OpenFile(s.path+"/"+id, os.O_RDONLY, 0777)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrMediaNotFound
		}
		return nil, err
	}

	defer file.Close()
	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	result := make([]byte, stat.Size())
	_, err = file.Read(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *localFileStorage) Upload(id string, data []byte) error {
	s.mutex.RLock()
	mutex := s.mutexes[id]
	s.mutex.RUnlock()

	if mutex == nil {
		mutex = &sync.Mutex{}
		s.mutex.Lock()
		s.mutexes[id] = mutex
		s.mutex.Unlock()
	}

	mutex.Lock()
	defer mutex.Unlock()
	file, err := os.OpenFile(s.path+"/"+id, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(data)

	if err != nil {
		return err
	}

	return nil
}

func (s *localFileStorage) Delete(id string) error {
	s.mutex.RLock()
	mutex := s.mutexes[id]
	s.mutex.RUnlock()

	if mutex == nil {
		mutex = &sync.Mutex{}
		s.mutex.Lock()
		s.mutexes[id] = mutex
		s.mutex.Unlock()
	}

	mutex.Lock()
	defer mutex.Unlock()
	err := os.Remove(s.path + "/" + id)

	if err != nil {
		if os.IsNotExist(err) {
			delete(s.mutexes, id)
			return domain.ErrMediaNotFound
		}
		return err
	}

	delete(s.mutexes, id)
	return nil
}

func initFileDir(path string) error {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			mkdirErr := os.Mkdir(path, 0777)
			if mkdirErr != nil {
				return mkdirErr
			}
			return nil
		}

		return err
	}

	return nil
}
