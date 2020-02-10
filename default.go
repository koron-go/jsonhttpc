package jsonhttpc

import "context"

// DefaultClient is used for `jsonhttpc.Do()` for convenience.
var DefaultClient = &Client{}

// Do makes a HTTP request to the server with JSON body, and parse a JSON in
// response body using `DefaultClient`.
func Do(ctx context.Context, method, pathOrURL string, body, receiver interface{}) error {
	return DefaultClient.Do(ctx, method, pathOrURL, body, receiver)
}
