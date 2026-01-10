package store_test

import (
	"errors"
	"testing"

	"github.com/sotiri-geo/dkv/store"
)

func TestGetKVStore(t *testing.T) {
	t.Run("Put key in store", func(t *testing.T) {
		// GIVEN store
		// GIVEN key k1 and value v1
		s := store.NewKVStore()
		k, v := "name", "john"
		// WHEN we put k-v pair k1:v1 into the store
		s.Put(k, v)

		// THEN we should have k-v pair k1:v1 persisted
		if v != s.Data[k] {
			t.Errorf("got %q, want %q", s.Data[k], v)
		}
	})

	t.Run("Get key in store", func(t *testing.T) {
		// GIVEN store with key k1 and value v1 already existing
		s := store.NewKVStore()
		k, v := "name", "john"
		s.Put(k, v)

		// WHEN we fetch for key k1
		got, err := s.Get(k)
		if err != nil {
			t.Fatalf("should not error: %v", err)
		}

		// THEN we should receive value v1
		if got != v {
			t.Errorf("got %q, want %q", got, v)
		}
	})

	t.Run("Get key not in store", func(t *testing.T) {
		// GIVEN empty store
		s := store.NewKVStore()

		// WHEN fetching for non-existent key k1
		_, err := s.Get("k1")

		// THEN should return with key not found error
		if err == nil {
			t.Fatalf("should error with key not found, got error: %v", err)
		}

		if !errors.Is(err, store.ErrKeyNotFound) {
			t.Errorf("got %+v, want %+v", err, store.ErrKeyNotFound)
		}
	})

	t.Run("Delete key in store", func(t *testing.T) {
		// GIVEN store with key k1 and value v1 already existing
		s := store.NewKVStore()
		k, v := "name", "john"
		s.Put(k, v)

		// WHEN we delete key k1
		err := s.Delete(k)
		if err != nil {
			t.Fatalf("should not error: %v", err)
		}

		// THEN k1:v1 pair should not be in store
		got, err := s.Get(k)

		if err == nil {
			t.Errorf("did not remove key %q, found with value %q", k, got)
		}
	})

	t.Run("Delete key not in store", func(t *testing.T) {
		// GIVEN empty store
		s := store.NewKVStore()

		// WHEN we delete non existent key k1
		err := s.Delete("k1")

		// THEN should return with key not found error
		if err == nil {
			t.Fatal("should error")
		}
		if !errors.Is(err, store.ErrKeyNotFound) {
			t.Errorf("got %+v, want %+v", err, store.ErrKeyNotFound)
		}
	})
}
