package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter"
)

func TestClient_Request(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		payload    interface{}
		result     interface{}
		wantErr    bool
		errContain string
	}{
		{
			name: "GET 200 with JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"key": "value"})
			},
			method:  http.MethodGet,
			path:    "/test",
			result:  &map[string]string{},
			wantErr: false,
		},
		{
			name: "POST 201 created",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]int{"id": 42})
			},
			method:  http.MethodPost,
			path:    "/items",
			payload: map[string]string{"name": "test"},
			result:  &map[string]int{},
			wantErr: false,
		},
		{
			name: "DELETE 204 no content, result nil",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			method:  http.MethodDelete,
			path:    "/items/1",
			result:  nil,
			wantErr: false,
		},
		{
			name: "400 bad request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			},
			method:     http.MethodPost,
			path:       "/items",
			payload:    map[string]string{"name": ""},
			wantErr:    true,
			errContain: "http 400",
		},
		{
			name: "401 unauthorized",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			},
			method:     http.MethodGet,
			path:       "/secret",
			wantErr:    true,
			errContain: "http 401",
		},
		{
			name: "404 not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "not found", http.StatusNotFound)
			},
			method:     http.MethodGet,
			path:       "/missing",
			wantErr:    true,
			errContain: "http 404",
		},
		{
			name: "429 too many requests",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "rate limited", http.StatusTooManyRequests)
			},
			method:     http.MethodGet,
			path:       "/test",
			wantErr:    true,
			errContain: "http 429",
		},
		{
			name: "500 internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal error", http.StatusInternalServerError)
			},
			method:     http.MethodGet,
			path:       "/test",
			wantErr:    true,
			errContain: "http 500",
		},
		{
			name: "invalid JSON in response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{broken json`))
			},
			method:     http.MethodGet,
			path:       "/test",
			result:     &map[string]string{},
			wantErr:    true,
			errContain: "decode response",
		},
		{
			name: "nil payload -> no body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				if len(body) != 0 {
					t.Errorf("expected empty body, got %d bytes", len(body))
				}
				w.WriteHeader(http.StatusOK)
			},
			method:  http.MethodGet,
			path:    "/test",
			payload: nil,
			result:  nil,
			wantErr: false,
		},
		{
			name: "POST with payload encodes body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var got map[string]string
				require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
				assert.Equal(t, "bar", got["foo"])
				w.WriteHeader(http.StatusOK)
			},
			method:  http.MethodPost,
			path:    "/test",
			payload: map[string]string{"foo": "bar"},
			result:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			c := NewClient(ts.URL)
			err := c.Request(context.Background(), tt.method, tt.path, tt.payload, tt.result)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClient_Request_SetsHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "TestAgent/1.0", r.Header.Get("User-Agent"))
		assert.Equal(t, "APIKey secret", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	signer := edgecenter.RequestSignerFunc(func(req *http.Request) error {
		req.Header.Set("Authorization", "APIKey secret")
		return nil
	})

	c := NewClient(ts.URL, WithUserAgent("TestAgent/1.0"), WithSigner(signer))
	err := c.Request(context.Background(), http.MethodGet, "/check", nil, nil)
	require.NoError(t, err)
}

func TestClient_Request_DecodesResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	var result map[string]string
	err := c.Request(context.Background(), http.MethodGet, "/test", nil, &result)
	require.NoError(t, err)
	assert.Equal(t, "true", result["ok"])
}
