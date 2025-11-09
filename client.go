package jsonhttpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

// Client is a JSON specialized HTTP client.
type Client struct {
	u *url.URL
	c *http.Client
	h http.Header

	lReq  RequestLogger
	lResp ResponseLogger
}

// New creates a new Client with base URL.
// You can pass nil to baseURL, in which case you should pass a complete URL to
// the pathOrURL argument of Do() function.
func New(baseURL *url.URL) *Client {
	return &Client{
		u: baseURL,
	}
}

// Clone create a clone of Client.
func (c *Client) Clone() *Client {
	nc := *c
	nc.h = nil
	if len(c.h) > 0 {
		nc.h = c.h.Clone()
	}
	return &nc
}

// WithBaseURL overrides base URL to use.
func (c *Client) WithBaseURL(u *url.URL) *Client {
	c.u = u
	return c
}

// WithClient overrides *http.Client to use.
func (c *Client) WithClient(hc *http.Client) *Client {
	c.c = hc
	return c
}

// WithHeader overrides http.Header to use.
func (c *Client) WithHeader(h http.Header) *Client {
	c.h = h
	return c
}

// Do makes a HTTP request to the server with JSON body, and parse a JSON in
// response body.
//
// `ctx` can be nil, in that case it use `context.Background()` instead.
//
// This returns an error when `pathOrURL` starts with "jsonhttpc.Parse error: "
// to treat `Path()`'s failure as error.
func (c *Client) Do(ctx context.Context, method, pathOrURL string, body, receiver any) error {
	// prepare the request.
	req, err := c.newRawRequest(ctx, method, pathOrURL, body)
	if err != nil {
		return err
	}

	// make messages round trip.
	c.logReq(ctx, req)
	resp, err := c.client().Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	c.logResp(ctx, resp)

	// verify status code.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.toError(resp)
	}
	if resp.ContentLength == 0 {
		if resp.StatusCode == 204 {
			return nil
		}
		return ErrEmptyResponse
	}

	// parse the response.
	if receiver == nil {
		return ErrReceiverAbsence
	}
	err = c.isJSONResponse(resp)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(receiver)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) newRawRequest(ctx context.Context, method, pathOrURL string, body any) (*http.Request, error) {
	u, err := c.parseURL(pathOrURL)
	if err != nil {
		return nil, err
	}
	r, err := c.bodyReader(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), r)
	if err != nil {
		return nil, err
	}
	c.setupHeader(req, body)
	return req, nil
}

func (c *Client) parseURL(pathOrURL string) (*url.URL, error) {
	if strings.HasPrefix(pathOrURL, "jsonhttpc.Parse error: ") {
		return nil, errors.New(pathOrURL)
	}
	if c.u != nil {
		return c.u.Parse(pathOrURL)
	}
	return url.Parse(pathOrURL)
}

func (c *Client) bodyReader(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	bb := &bytes.Buffer{}
	err := json.NewEncoder(bb).Encode(body)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func (c *Client) setupHeader(req *http.Request, body any) {
	if len(c.h) > 0 {
		req.Header = c.h.Clone()
	}
	if body != nil {
		req.Header.Set("Content-Type", c.getContentType(body))
	}
}

func (c *Client) client() *http.Client {
	if c.c != nil {
		return c.c
	}
	return http.DefaultClient
}

type contentTyper interface {
	ContentType() string
}

func (c *Client) getContentType(v any) string {
	if w, ok := v.(contentTyper); ok {
		return w.ContentType()
	}
	return "application/json"
}

func (c *Client) isJSONResponse(r *http.Response) error {
	ct := r.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return newSystemError(r, fmt.Errorf("invalid content type: %w", err))
	}
	if !c.isApplicationJSON(mt) {
		return newSystemError(r, fmt.Errorf("not JSON response: %s", ct))
	}
	return nil
}

func (c *Client) isApplicationJSON(ct string) bool {
	if ct == "application/json" {
		return true
	}
	if strings.HasPrefix(ct, "application/") && strings.HasSuffix(ct, "+json") {
		return true
	}
	return false
}

func (c *Client) toError(r *http.Response) error {
	err := c.isJSONResponse(r)
	if err != nil {
		return err
	}
	er, err := parseError(r)
	if err != nil {
		return err
	}
	return er
}

func (c *Client) logReq(ctx context.Context, req *http.Request) {
	if c.lReq == nil {
		return
	}
	c.lReq.LogRequest(ctx, req)
}

func (c *Client) logResp(ctx context.Context, resp *http.Response) {
	if c.lResp == nil {
		return
	}
	c.lResp.LogResponse(ctx, resp)
}
