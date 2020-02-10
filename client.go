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
func (c *Client) Do(ctx context.Context, method, pathOrURL string, body, receiver interface{}) error {
	if strings.HasPrefix(pathOrURL, "jsonhttpc.Parse error: ") {
		return errors.New(pathOrURL)
	}
	// prepare the request.
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := c.parseURL(pathOrURL)
	if err != nil {
		return err
	}
	r, err := c.bodyReader(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), r)
	if err != nil {
		return err
	}
	c.setupHeader(req, body)

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

func (c *Client) parseURL(pathOrURL string) (*url.URL, error) {
	if c.u != nil {
		return c.u.Parse(pathOrURL)
	}
	return url.Parse(pathOrURL)
}

func (c *Client) bodyReader(body interface{}) (io.Reader, error) {
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

func (c *Client) setupHeader(req *http.Request, body interface{}) {
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

func (c *Client) getContentType(v interface{}) string {
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
