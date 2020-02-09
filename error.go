package jsonhttpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	// ErrEmptyResponse is received empty response unexpectedly.
	// Ex. StatusCode shows success but Content-Length is 0.
	// This is a server side error.
	ErrEmptyResponse = errors.New("unexpected empty response")

	// ErrReceiverAbsence means `receiver` is nil.
	// Server responded with body, but user did't provide any receivers.
	// This is a user side error.
	ErrReceiverAbsence = errors.New("receiver absence")
)

// SystemError is for problems on jsonhttp.
type SystemError struct {
	StatusCode int
	Status     string
	Err        error
}

var _ error = (*SystemError)(nil)

func newSystemError(r *http.Response, err error) error {
	return &SystemError{
		StatusCode: r.StatusCode,
		Status:     r.Status,
		Err:        err,
	}
}

func (se *SystemError) Error() string {
	return fmt.Sprintf("jsonhttpc system problem: %s: %s", se.Status, se.Err)
}

// Unwrap obtains the based error.
func (se *SystemError) Unwrap() error {
	return se.Err
}

// Error is general structure to store error.
// This supports https://tools.ietf.org/html/rfc7807, if "status" field has
// some troubles (not string or differ from `StatusCode`, its raw value is put
// into `Properties`.
type Error struct {
	StatusCode  int
	Status      string
	ContentType string

	Type     string
	Title    string
	Detail   string
	Instance string

	Properties map[string]interface{}
}

var _ error = (*Error)(nil)

func parseError(r *http.Response) (*Error, error) {
	er := &Error{
		StatusCode:  r.StatusCode,
		Status:      r.Status,
		ContentType: r.Header.Get("Content-Type"),
	}

	var props map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&props)
	if err != nil {
		return nil, newSystemError(r, fmt.Errorf("failed to decode error JSON: %w", err))
	}
	for k, v := range props {
		var del bool
		switch k {
		case "type":
			er.Type = fmt.Sprint(v)
			del = true
		case "title":
			er.Title = fmt.Sprint(v)
			del = true
		case "status":
			n, err := strconv.Atoi(fmt.Sprint(v))
			if err != nil {
				continue
			}
			if er.StatusCode != n {
				continue
			}
			del = true
		case "detail":
			er.Detail = fmt.Sprint(v)
			del = true
		case "instance":
			er.Instance = fmt.Sprint(v)
			del = true
		}
		if del {
			delete(props, k)
		}
	}
	er.Properties = props

	return er, nil
}

func (er *Error) Error() string {
	return fmt.Sprintf("error: status=%q props=%#v", er.Status, er.Properties)
}
