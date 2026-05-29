package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter/provider"
)

const apiKey = "integration-test-key"

type recordedRequest struct {
	Method string
	Path   string
	Auth   string
	Body   []byte
}

type fakeRMON struct {
	mu       sync.Mutex
	seq      int
	store    map[string]map[string]interface{}
	requests []recordedRequest
}

func newFakeRMON() *fakeRMON {
	return &fakeRMON{store: make(map[string]map[string]interface{})}
}

func (f *fakeRMON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	f.mu.Lock()
	defer f.mu.Unlock()

	f.requests = append(f.requests, recordedRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Auth:   r.Header.Get("Authorization"),
		Body:   body,
	})

	switch r.Method {
	case http.MethodPost:
		f.seq++
		id := f.seq
		obj := decode(body)
		obj["id"] = id
		f.store[r.URL.Path+"/"+strconv.Itoa(id)] = obj

		writeJSON(w, http.StatusCreated, obj)
	case http.MethodGet:
		obj, ok := f.store[r.URL.Path]
		if !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, normalizeForRead(obj))
	case http.MethodPut:
		obj := decode(body)
		obj["id"] = idFromPath(r.URL.Path)
		f.store[r.URL.Path] = obj

		writeJSON(w, http.StatusOK, obj)
	case http.MethodDelete:
		delete(f.store, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func start(t *testing.T) (*provider.Client, *fakeRMON) {
	t.Helper()

	srv := newFakeRMON()
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	signer := provider.WithSignerFunc(func(req *http.Request) error {
		for k, v := range provider.AuthenticatedHeaders(apiKey) {
			req.Header.Set(k, v)
		}
		return nil
	})

	client := provider.NewClient(ts.URL, signer, provider.WithUserAgent("integration-test/1.0"))

	return client, srv
}

func (f *fakeRMON) calls() []recordedRequest {
	f.mu.Lock()
	defer f.mu.Unlock()

	out := make([]recordedRequest, len(f.requests))
	copy(out, f.requests)

	return out
}

func normalizeForRead(obj map[string]interface{}) map[string]interface{} {
	raw, ok := obj["checks"].([]interface{})
	if !ok {
		return obj
	}

	out := make(map[string]interface{}, len(obj))
	for k, v := range obj {
		out[k] = v
	}

	checks := make([]interface{}, 0, len(raw))
	for _, id := range raw {
		checks = append(checks, map[string]interface{}{"multi_check_id": id})
	}
	out["checks"] = checks

	return out
}

func decode(body []byte) map[string]interface{} {
	m := make(map[string]interface{})
	if len(body) > 0 {
		_ = json.Unmarshal(body, &m)
	}

	return m
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func idFromPath(path string) int {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id, _ := strconv.Atoi(parts[len(parts)-1])

	return id
}

type call struct {
	method, path string
}

func assertPaths(t *testing.T, got []recordedRequest, want []call) {
	t.Helper()

	require.Len(t, got, len(want))

	for i, w := range want {
		assert.Equal(t, w.method, got[i].Method, "call %d", i)
		assert.Equal(t, w.path, got[i].Path, "call %d", i)
	}
}

func itoa(n int) string { return strconv.Itoa(n) }
