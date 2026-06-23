package api

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func newRequest(url string) *HttpRequest {
	return &HttpRequest{
		Method: "GET",
		URL:    url,
	}
}

func TestExecute_Success(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	defer server.Close()

	req := newRequest(server.URL)

	resp, err := Execute(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200")
	}

	if string(resp.Body) != "ok" {
		t.Fatalf("unexpected body")
	}
}

func TestExecute_Retry(t *testing.T) {
	var count atomic.Int32

	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		c := count.Add(1)

		if c < 3 {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	defer server.Close()

	req := newRequest(server.URL)

	req.RetryDelayFunction = func(statusCode int, err error, attempt int) *time.Duration {
		if statusCode == 500 && attempt < 3 {
			d := 10 * time.Millisecond
			return &d
		}
		return nil
	}

	resp, err := Execute(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected success after retry")
	}

	if count.Load() < 3 {
		t.Fatalf("retry did not happen")
	}
}

func TestExecute_Timeout(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
	})
	defer server.Close()

	timeout := int64(50)

	req := &HttpRequest{
		Method:    "GET",
		URL:       server.URL,
		TimeoutMs: &timeout,
	}

	_, err := Execute(req)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
}

func TestExecuteAsync(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	defer server.Close()

	req := newRequest(server.URL)

	respCh, errCh := ExecuteAsync(req)

	select {
	case resp := <-respCh:
		if resp.StatusCode != 200 {
			t.Fatalf("expected 200")
		}
	case err := <-errCh:
		t.Fatal(err)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for async response")
	}
}
