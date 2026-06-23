package api

import (
	"testing"
)

func TestHttpResponse_Basic(t *testing.T) {
	resp := HttpResponse{
		StatusCode: 200,
		Headers:    map[string][]string{"X-Test": {"1"}},
		Body:       []byte("ok"),
	}

	if resp.StatusCode != 200 {
		t.Fatal("invalid status")
	}

	if resp.Headers["X-Test"][0] != "1" {
		t.Fatal("invalid header")
	}

	if string(resp.Body) != "ok" {
		t.Fatal("invalid body")
	}
}
