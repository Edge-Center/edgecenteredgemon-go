package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter"
)

type Client struct {
	httpc   *http.Client
	signer  edgecenter.RequestSigner
	ua      string
	baseURL string
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	httpc := &http.Client{Timeout: time.Minute}

	c := &Client{
		httpc:   httpc,
		baseURL: baseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Request(ctx context.Context, method, path string, payload interface{}, result interface{}) error {
	var body io.Reader

	if payload != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return fmt.Errorf("encode payload: %w", err)
		}
		body = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(raw))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")

	if c.ua != "" {
		req.Header.Set("User-Agent", c.ua)
	}

	if c.signer != nil {
		if err := c.signer.Sign(req); err != nil {
			return nil, fmt.Errorf("sign request: %w", err)
		}
	}

	return c.httpc.Do(req)
}
