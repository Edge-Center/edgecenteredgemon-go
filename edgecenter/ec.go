package edgecenter

import (
	"context"
	"net/http"
)

// Requester is the abstraction for sending a single HTTP request to the API.
// All SDK services (channel, checkgroup, checkhttp, etc.) call only Request:
// they pass method, path, payload, and a pointer to result; the Requester implementation
// decides how to build the request, sign it, and decode the response.
// Implemented by provider.Client (Request method in edgecenter/provider/provider.go).
// Example call from a service: s.r.Request(ctx, http.MethodPost, "/rmon/channel/telegram", req, &resp).
type Requester interface {
	Request(ctx context.Context, method, path string, payload interface{}, result interface{}) error
}

// RequestSigner adds authorization data to a request (e.g. the Authorization header).
// It keeps signing logic in one place instead of in every service.
// In the provider it is used in Client.do(req): before c.httpc.Do(req), c.signer.Sign(req) is called;
// if the signer was set via WithSigner, it adds the required headers (APIKey, Bearer, etc.) to req.
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
