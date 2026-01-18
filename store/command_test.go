package store_test

import (
	"errors"
	"testing"

	"github.com/sotiri-geo/dkv/store"
)

func TestGetQuery(t *testing.T) {
	t.Run("GET query on existing key", func(t *testing.T) {
		// GIVEN store s1 with key k1, value v1
		k, v := "k1", "v1"
		s := store.NewKVStore()
		putCommand := store.NewPutCommand(k, v)
		s.Apply(putCommand)

		getQuery := store.NewGetQuery(k)

		// WHEN fetching key k1
		got, err := s.Execute(getQuery)
		// THEN we recieve value v1
		if err != nil {
			t.Fatalf("should not error: %v", err)
		}

		if got != v {
			t.Errorf("got %q, want %q", got, v)
		}
	})

	t.Run("GET query on non-existent key", func(t *testing.T) {
		// GIVEN empty store s1
		s := store.NewKVStore()
		getQuery := store.NewGetQuery("non-existent")

		// WHEN fetching non-existent key
		_, err := s.Execute(getQuery)
		// THEN we recieve ErrKeyNotFound
		if err == nil {
			t.Fatalf("should not error: %v", err)
		}

		if !errors.Is(err, store.ErrKeyNotFound) {
			t.Errorf("incorrect error response: got %v", err)
		}
	})
}

func TestPutCommand(t *testing.T) {
	t.Run("apply PUT command", func(t *testing.T) {
		// GIVEN store s1 and put command p1 with key k1, value v1
		k, v := "foo", "bar"
		s := store.NewKVStore()
		putCommand := store.NewPutCommand(k, v)

		// WHEN applying p1 to store s1
		s.Apply(putCommand)

		// THEN k-v pair k1:v1 is persisted to store s1
		getQuery := store.NewGetQuery(k)
		got, err := s.Execute(getQuery)
		if err != nil {
			t.Fatalf("should not error: %v", err)
		}

		if got != v {
			t.Errorf("got %q, want %q", got, v)
		}
	})
}

func TestDeleteCommand(t *testing.T) {
	t.Run("delete existing key command", func(t *testing.T) {
		// GIVEN store s1 with k-v pairs k1, v1 and delete command d1
		k, v := "foo", "bar"
		s := store.NewKVStore()
		putCommand := store.NewPutCommand(k, v)
		s.Apply(putCommand)

		// WHEN deleting key k1 from store s1
		deleteCommand := store.NewDeleteCommand(k)
		s.Apply(deleteCommand)

		// THEN key k1 no longer exists
		getQuery := store.NewGetQuery(k)
		_, err := s.Execute(getQuery)

		if !errors.Is(err, store.ErrKeyNotFound) {
			t.Errorf("incorrect error response: got %v", err)
		}
	})
}
