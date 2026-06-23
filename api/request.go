package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type HttpRequest struct {
	Method             string                                                      `json:"method"`
	URL                string                                                      `json:"url"`
	Headers            map[string]string                                           `json:"headers,omitempty"`
	Query              map[string]string                                           `json:"query,omitempty"`
	Body               any                                                         `json:"body,omitempty"`
	TimeoutMs          *int64                                                      `json:"timeout_ms,omitempty"`
	RetryDelayFunction func(statusCode int, err error, attempt int) *time.Duration `json:"-"`
}

// Safely get retry delay
func (r *HttpRequest) GetRetryDelay(statusCode int, err error, attempt int) *time.Duration {
	if r.RetryDelayFunction == nil {
		return nil
	}
	return r.RetryDelayFunction(statusCode, err, attempt)
}

// Build final URL with query params
func (r *HttpRequest) BuildURL() (string, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return "", err
	}

	q := u.Query()

	// Override or set query params (1:1 mapping)
	for key, value := range r.Query {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (r *HttpRequest) BuildBody() ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	switch v := r.Body.(type) {

	case []byte:
		return v, nil

	case string:
		return []byte(v), nil

	case *bytes.Buffer:
		return v.Bytes(), nil

	case bytes.Buffer:
		return v.Bytes(), nil

	default:
		// fallback → JSON encode
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}

		// auto-set content-type if not already set
		if r.Headers == nil {
			r.Headers = map[string]string{}
		}
		if _, ok := r.Headers["Content-Type"]; !ok {
			r.Headers["Content-Type"] = "application/json"
		}

		return b, nil
	}
}
