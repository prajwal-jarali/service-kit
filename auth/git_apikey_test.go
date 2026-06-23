package auth

import (
	"context"
	"testing"
)

func TestNewGitApiKey_Validation(t *testing.T) {
	_, err := NewGitApiKey("")
	if err == nil {
		t.Fatalf("expected error for empty token")
	}
}

func TestNewGitApiKey_Success(t *testing.T) {
	g, err := NewGitApiKey("ghp_test")
	if err != nil {
		t.Fatal(err)
	}

	if g.token != "Bearer ghp_test" {
		t.Fatalf("unexpected token format")
	}
}

func TestGitApiKey_GetToken(t *testing.T) {
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
