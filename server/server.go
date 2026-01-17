// Package server implements the handlers for GET, PUT, DELETE operations
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sotiri-geo/dkv/store"
)

type Server struct {
	store *store.KVStore
}

type GetResponse struct {
	Value string `json:"value"`
}

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewServer(store *store.KVStore) *Server {
	return &Server{
		store,
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.Routes())
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Adding handlers
	mux.HandleFunc("GET /kv/{key}", s.handleGet)
	mux.HandleFunc("PUT /kv", s.handlePut)
	mux.HandleFunc("DELETE /kv/{key}", s.handleDelete)
	return mux
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "please provide a key to search", http.StatusBadRequest)
		return
	}

	value, err := s.store.Get(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("key %q not found", key), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&GetResponse{value}); err != nil {
		// Log error for now rather than sending it back to the client
		log.Printf("error enconding response: %v", err)
	}
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	var req PutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json request", http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		http.Error(w, "key cannot be empty", http.StatusBadRequest)
		return
	}

	s.store.Put(req.Key, req.Value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	if key == "" {
		http.Error(w, "key cannot be empty", http.StatusBadRequest)
		return
	}

	err := s.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
