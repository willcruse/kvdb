package storagebackend

import (
	"testing"
)

const (
	TEST_KEY   = "test_key"
	TEST_VALUE = "test_value"
)

func TestMapStorageInit(t *testing.T) {
	mapStorageBackend := &MapStorageBackend{}
	mapStorageBackend.Init()

	if mapStorageBackend.data == nil {
		t.Error("mapStorageBackend.data == nil. Expected map to be initialised")
	}
}

func TestMapStorageCanSet(t *testing.T) {
	mapStorageBackend := &MapStorageBackend{}
	mapStorageBackend.Init()

	err := mapStorageBackend.Set(TEST_KEY, TEST_VALUE)
	if err != nil {
		t.Errorf("Failed to set key. Got err = %s", err)
	}

	fetchedValue, err := mapStorageBackend.Get(TEST_KEY)

	if err != nil {
		t.Errorf("Failed to get key. Got err = %s", err)
	}

	if fetchedValue != TEST_VALUE {
		t.Errorf("Expected fetched value to match expected value. %s != %s", TEST_VALUE, fetchedValue)
	}

}

// This is a duplicate of the above test
// Not the most helpful...
func TestMapStorageCanGet(t *testing.T) {
	mapStorageBackend := &MapStorageBackend{}
	mapStorageBackend.Init()

	err := mapStorageBackend.Set(TEST_KEY, TEST_VALUE)
	if err != nil {
		t.Errorf("Failed to set key. Got err = %s", err)
	}

	fetchedValue, err := mapStorageBackend.Get(TEST_KEY)

	if err != nil {
		t.Errorf("Failed to get key. Got err = %s", err)
	}

	if fetchedValue != TEST_VALUE {
		t.Errorf("Expected fetched value to match expected value. %s != %s", TEST_VALUE, fetchedValue)
	}

}

func TestMapStorageCanDelete(t *testing.T) {
	mapStorageBackend := &MapStorageBackend{}
	mapStorageBackend.Init()

	err := mapStorageBackend.Set(TEST_KEY, TEST_VALUE)
	if err != nil {
		t.Errorf("Failed to set key. Got err = %s", err)
	}

	fetchedValue, err := mapStorageBackend.Get(TEST_KEY)

	if err != nil {
		t.Errorf("Failed to get key. Got err = %s", err)
	}

	err = mapStorageBackend.Delete(TEST_KEY)
	if err != nil {
		t.Errorf("Failed to delete key. Got err = %s", err)
	}

	fetchedValue, err = mapStorageBackend.Get(TEST_KEY)
	if fetchedValue != "" || err == nil {
		t.Errorf("Expected fetching deleted value to result in an error. Got value %s and err = %s", fetchedValue, err)
	}
}
