package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{}

func Execute(req *HttpRequest) (*HttpResponse, error) {
	attempt := 0

	for {
		// Build URL
		finalURL, err := req.BuildURL()
		if err != nil {
			return nil, err
		}

		// Build Body
		bodyBytes, err := req.BuildBody()
		if err != nil {
			return nil, err
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		// Create HTTP request
		httpReq, err := http.NewRequest(req.Method, finalURL, bodyReader)
		if err != nil {
			return nil, err
		}

		// Apply headers
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		// Apply timeout via context
		ctx := httpReq.Context()
		if req.TimeoutMs != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(*req.TimeoutMs)*time.Millisecond)
			defer cancel()
		}
		httpReq = httpReq.WithContext(ctx)

		// Execute request
		resp, err := httpClient.Do(httpReq)

		var statusCode int
		var respBody []byte
		var headers map[string][]string

		if resp != nil {
			statusCode = resp.StatusCode
			headers = resp.Header

			respBody, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
		}

		// Retry decision
		delay := req.GetRetryDelay(statusCode, err, attempt)
		if delay == nil {
			if err != nil {
				return nil, err
			}
			return &HttpResponse{
				StatusCode: statusCode,
				Headers:    headers,
				Body:       respBody,
			}, nil
		}

		time.Sleep(*delay)
		attempt++
	}
}

func ExecuteAsync(req *HttpRequest) (<-chan *HttpResponse, <-chan error) {
	respCh := make(chan *HttpResponse, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(respCh)
		defer close(errCh)

		resp, err := Execute(req)
		if err != nil {
			errCh <- err
			return
		}

		respCh <- resp
	}()

	return respCh, errCh
}
