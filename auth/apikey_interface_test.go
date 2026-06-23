package auth

import (
	"context"
	"testing"
)

// compile-time check helpers
func assertTokenProvider(t TokenProvider) {}

func TestTokenProvider_Implementation(t *testing.T) {
	// ApiKey should implement TokenProvider
	a, err := NewApiKey("dummy", false, false)
	if err != nil {
		t.Fatal(err)
	}
	assertTokenProvider(a)

	// GitApiKey should implement TokenProvider
	g, err := NewGitApiKey("ghp_dummy")
	if err != nil {
		t.Fatal(err)
	}
	assertTokenProvider(g)
}

func TestTokenProvider_GetToken(t *testing.T) {
	g, err := NewGitApiKey("ghp_test")
	if err != nil {
		t.Fatal(err)
	}

	token, err := g.GetToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	expected := "Bearer ghp_test"
	if token != expected {
		t.Fatalf("expected %s, got %s", expected, token)
	}
}