// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testacc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
)

// ----------------------------------------------------------------------------
// In-memory model (very small subset of your real Ad-Server object)
// ----------------------------------------------------------------------------

type adServer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ----------------------------------------------------------------------------
// Mock server
// ----------------------------------------------------------------------------

type Mock struct {
	srv       *httptest.Server
	mu        sync.RWMutex
	adServers map[uint]adServer
	nextID    uint
}

// NewMock spins up the fake Broadpeak IO backend.
func NewMock() *Mock {
	m := &Mock{
		adServers: make(map[uint]adServer),
		nextID:    1,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/adservers", m.handleCollection)
	mux.HandleFunc("/v1/adservers/", m.handleItem) // trailing slash for /{id}

	m.srv = httptest.NewServer(mux)
	return m
}

// URL returns the root endpoint you inject into the provider (BPKIO_ENDPOINT).
func (m *Mock) URL() string { return m.srv.URL }

// Close must be deferred in the test.
func (m *Mock) Close() { m.srv.Close() }

// -----------------------------------------------------------------------------
// Public helper for the acceptance test
// -----------------------------------------------------------------------------

// ResourceExists tells the test whether an Ad-Server identified by string id
// (Terraform keeps IDs as strings) is still present in the mock DB.
func (m *Mock) ResourceExists(id string) bool {
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return false
	}

	m.mu.RLock()
	_, ok := m.adServers[uint(uid)]
	m.mu.RUnlock()
	return ok
}

// -----------------------------------------------------------------------------
// HTTP handlers
// -----------------------------------------------------------------------------

// POST /v1/adservers         -> create
// GET  /v1/adservers         -> list (unused, but handy)
// ---------------------------------------------------------------------------

func (m *Mock) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var in struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m.mu.Lock()
		id := m.nextID
		m.nextID++
		obj := adServer{ID: id, Name: in.Name}
		m.adServers[id] = obj
		m.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(obj)

	case http.MethodGet:
		m.mu.RLock()
		defer m.mu.RUnlock()

		var list []adServer
		for _, v := range m.adServers {
			list = append(list, v)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /v1/adservers/{id}     -> read
// PUT /v1/adservers/{id}     -> update
// DELETE /v1/adservers/{id}  -> delete
// ---------------------------------------------------------------------------

func (m *Mock) handleItem(w http.ResponseWriter, r *http.Request) {
	// Trim prefix and parse ID.
	idStr := r.URL.Path[len("/v1/adservers/"):]
	uid64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	uid := uint(uid64)

	switch r.Method {
	case http.MethodGet:
		m.mu.RLock()
		obj, ok := m.adServers[uid]
		m.mu.RUnlock()
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(obj)

	case http.MethodDelete:
		m.mu.Lock()
		delete(m.adServers, uid)
		m.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
