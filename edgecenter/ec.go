package edgecenter

import (
	"context"
	"net/http"
)

// Requester is the abstraction for sending a single HTTP request to the API.
// All SDK services (channel, checkgroup, checkhttp, etc.) call only Request:
// they pass method, path, payload, and a pointer to result; the Requester implementation
// decides how to build the request, sign it, and decode the response.
// The primary implementation is [provider.Client], which handles
// JSON encoding, request signing, and error parsing.
// Example call from a service: s.r.Request(ctx, http.MethodPost, "/rmon/channel/telegram", req, &resp).
type Requester interface {
	Request(ctx context.Context, method, path string, payload interface{}, result interface{}) error
}

// RequestSigner adds authorization data to a request (e.g. the Authorization header).
// It keeps signing logic in one place instead of in every service.
// Pass an implementation via [provider.WithSigner] when constructing the client.
type RequestSigner interface {
	Sign(req *http.Request) error
}

// RequestSignerFunc is an adapter: a plain function with signature func(*http.Request) error
// implements the RequestSigner interface. Same idea as http.HandlerFunc in the standard library:
// no need to define a separate type with a Sign method, just pass a function.
type RequestSignerFunc func(req *http.Request) error

func (f RequestSignerFunc) Sign(req *http.Request) error {
	return f(req)
}
