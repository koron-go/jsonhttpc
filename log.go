package jsonhttpc

import (
	"context"
	"net/http"
)

// RequestLogger is a callback with request before make round trip.
type RequestLogger interface {
	LogRequest(context.Context, *http.Request)
}

// ResponseLogger is a callback with response after round trip succeeded.
type ResponseLogger interface {
	LogResponse(context.Context, *http.Response)
}

// RequestResponseLogger provides RequestLogger and Response by one object.
// It is supposed to use with WithRequestResponseLogger.
type RequestResponseLogger interface {
	RequestLogger
	ResponseLogger
}

// WithRequestLogger overrides RequestLogger of the client.
func (c *Client) WithRequestLogger(l RequestLogger) *Client {
	c.lReq = l
	return c
}

// WithResponseLogger overrides ResponseLogger of the client.
func (c *Client) WithResponseLogger(l ResponseLogger) *Client {
	c.lResp = l
	return c
}

// WithRequestResponseLogger overrides RequestLogger and ResponseLogger with an
// object in same time.
func (c *Client) WithRequestResponseLogger(l RequestResponseLogger) *Client {
	c.lReq = l
	c.lResp = l
	return c
}
