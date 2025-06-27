package storagebackend

import "fmt"

type StorageBackend interface {
	Init()
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type MapStorageBackend struct {
	data map[string]string
}

func (msb *MapStorageBackend) Init() {
	msb.data = make(map[string]string)
}

func (msb *MapStorageBackend) Set(key, value string) error {
	msb.data[key] = value
	return nil
}

func (msb *MapStorageBackend) Get(key string) (string, error) {
	value, exists := msb.data[key]

	if exists {
		return value, nil
	}

	return "", fmt.Errorf("map_storage_backend: no such key %s", key)
}

func (msb *MapStorageBackend) Delete(key string) error {
	_, exists := msb.data[key]
	if exists {
		delete(msb.data, key)
	}

	return nil
}
