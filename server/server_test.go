package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sotiri-geo/dkv/server"
	"github.com/sotiri-geo/dkv/store"
)

func TestRouterGet(t *testing.T) {
	// GIVEN a store with key k1 and value v1
	kvStore := store.NewKVStore()
	want := server.GetResponse{
		Value: "v1",
	}
	k, v := "k1", "v1"
	kvStore.Put(k, v)

	srv := server.NewServer(kvStore)
	tsrv := httptest.NewServer(srv.Routes())
	defer tsrv.Close()

	// WHEN making a GET http request
	resp, err := http.Get(tsrv.URL + "/kv?key=k1")
	if err != nil {
		t.Fatalf("should not fail: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// THEN response should contain value v1 with status OK
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var got server.GetResponse

	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatalf("should not error: %v", err)
	}

	if got.Value != want.Value {
		t.Errorf("got %q, want %q", got.Value, want.Value)
	}
}

func TestRouterPut(t *testing.T) {
	// GIVEN empty store
	// GIVEN key k1 with value v1
	kvStore := store.NewKVStore()
	requestBody := server.PutRequest{
		Key:   "k1",
		Value: "v1",
	}
	svr := server.NewServer(kvStore)
	tsvr := httptest.NewServer(svr.Routes())
	defer tsvr.Close()

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("marshalling request body: %v", err)
	}

	// WHEN making a PUT request with key k1 and value v1
	req, err := http.NewRequest(http.MethodPut, tsvr.URL+"/kv", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("creating request body: %v", err)
	}
	req.Header.Set("Content-type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("making put request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// THEN should return 201 created and k-v pair persisted in store
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	value, err := kvStore.Get("k1")
	if err != nil {
		t.Fatalf("fetching stored key %q: %v", requestBody.Key, err)
	}

	if value != requestBody.Value {
		t.Errorf("got value %q, want %q", value, requestBody.Value)
	}
}

func TestRouterDelete(t *testing.T) {
	// GIVEN store with key k1 and value v1
	kvStore := store.NewKVStore()
	k, v := "k1", "v1"
	kvStore.Put(k, v)

	svr := server.NewServer(kvStore)
	tsvr := httptest.NewServer(svr.Routes())
	defer tsvr.Close()

	// WHEN a DELETE request is made for key k1
	req, err := http.NewRequest(http.MethodDelete, tsvr.URL+"/kv/k1", nil)
	if err != nil {
		t.Fatalf("creating DELETE request: %v", err)
	}

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("sending DELETE request: %v", err)
	}

	// THEN response status code 204 with k1 no longer existing in store
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = kvStore.Get(k)

	if err == nil {
		t.Errorf("key %q still exists in store", k)
	}
}
