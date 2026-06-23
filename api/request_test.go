package api

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestBuildURL(t *testing.T) {
	req := HttpRequest{
		URL: "https://example.com/api?existing=1",
		Query: map[string]string{
			"q": "test",
		},
	}

	u, err := req.BuildURL()
	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := url.Parse(u)
	q := parsed.Query()

	if q.Get("existing") != "1" {
		t.Fatalf("expected existing=1")
	}

	if q.Get("q") != "test" {
		t.Fatalf("expected q=test")
	}
}

func TestBuildBody_JSON(t *testing.T) {
	req := HttpRequest{
		Body: map[string]string{"name": "test"},
	}

	b, err := req.BuildBody()
	if err != nil {
		t.Fatal(err)
	}

	var out map[string]string
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}

	if out["name"] != "test" {
		t.Fatalf("unexpected body")
	}
}

func TestBuildBody_String(t *testing.T) {
	req := HttpRequest{
		Body: "hello",
	}

	b, err := req.BuildBody()
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "hello" {
		t.Fatalf("expected raw string body")
	}
}
