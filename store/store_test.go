package store_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/sotiri-geo/dkv/store"
)

func TestStoreLog(t *testing.T) {
	t.Run("logs concurrent put, delete commands", func(t *testing.T) {
		// GIVEN empty store
		s := store.NewKVStore()
		n := 100
		var wg sync.WaitGroup

		// WHEN sending 100 concurrent put/delete requests
		for i := range n {
			idx := i
			wg.Go(func() {
				k, v := fmt.Sprintf("key%d", idx), fmt.Sprintf("value%d", idx)
				putCmd := store.NewPutCommand(k, v)
				delCmd := store.NewDeleteCommand(k)
				s.Apply(putCmd)
				s.Apply(delCmd)
			})
		}
		wg.Wait()

		// THEN we have 200 command logs
		got := len(s.GetLog())
		want := 200
		if got != want {
			t.Errorf("got log count %d, want %d", got, want)
		}
	})

	t.Run("query logs are not recorded", func(t *testing.T) {
		// GIVEN store with key k1, value v1
		s := store.NewKVStore()
		k, v := "k1", "v1"
		putCmd := store.NewPutCommand(k, v)
		s.Apply(putCmd)
		getQuery := store.NewGetQuery(k)

		// WHEN executing get query
		_, err := s.Execute(getQuery)
		if err != nil {
			t.Fatalf("should not error: %v", err)
		}

		// THEN query log is not recorded
		got := len(s.GetLog())
		want := 1 // putCmd counts as a single log
		if got != want {
			t.Errorf("got log count %d, want %d", got, want)
		}
	})
}

func TestReplay(t *testing.T) {
	// GIVEN empty store
	// GIVEN put k1-v1, put k2-v2 and delete k1 commands
	s := store.NewKVStore()

	// Build up state
	s.Apply(store.NewPutCommand("k1", "v1"))
	s.Apply(store.NewPutCommand("k2", "v2"))
	s.Apply(store.NewDeleteCommand("k1"))

	// Replay Cmds
	s.Replay()

	// THEN only k2-v2 remains in store
	_, err := s.Execute(store.NewGetQuery("k1"))

	if err == nil {
		t.Error("key k1 should not be in store")
	}

	value, err := s.Execute(store.NewGetQuery("k2"))
	if err != nil {
		t.Fatal("key k2 should exist in store")
	}

	if value != "v2" {
		t.Errorf("got value %q, want %q", value, "v2")
	}
}

func TestDeterminsm(t *testing.T) {
	// GIVEN empty store s1
	// GIVEN empty store s2
	// GIVEN commands put k1-v1, put k2-v2, delete v1
	s1 := store.NewKVStore()
	s2 := store.NewKVStore()

	cmds := []store.Command{
		store.NewPutCommand("k1", "v1"),
		store.NewPutCommand("k2", "v2"),
		store.NewDeleteCommand("k1"),
	}

	// WHEN applying commands on s1, s2
	for _, cmd := range cmds {
		s1.Apply(cmd)
		s2.Apply(cmd)
	}

	// THEN s1 and s2 have equal state
	_, err1 := s1.Execute(store.NewGetQuery("k1"))
	_, err2 := s2.Execute(store.NewGetQuery("k1"))

	if !errors.Is(err1, store.ErrKeyNotFound) && !errors.Is(err2, store.ErrKeyNotFound) {
		t.Errorf("stores have diverged: err1 %v, err2 %v", err1, err2)
	}

	v1, err1 := s1.Execute(store.NewGetQuery("k2"))
	v2, err2 := s2.Execute(store.NewGetQuery("k2"))

	if err1 != err2 || v1 != v2 {
		t.Errorf("stores have diverged: v1 %q err1 %v, v2 %q err2 %v", v1, err1, v2, err2)
	}
}
