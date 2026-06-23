package auth

import (
	"context"
	"testing"
	"time"
)

func TestNewApiKey_Validation(t *testing.T) {
	_, err := NewApiKey("", false, false)
	if err == nil {
		t.Fatalf("expected error for empty api key")
	}
}

func TestNewApiKey_Success(t *testing.T) {
	a, err := NewApiKey("test-key", true, true)
	if err != nil {
		t.Fatal(err)
	}

	if a.Apikey != "test-key" {
		t.Fatalf("unexpected api key")
	}

	if !a.testEndpoint || !a.privateEndpoint {
		t.Fatalf("flags not set correctly")
	}
}

func TestIamApiKey_CachedToken(t *testing.T) {
	a, err := NewApiKey("dummy", false, false)
	if err != nil {
		t.Fatal(err)
	}

	// inject cached token
	a.token = "Bearer cached-token"
	a.generatedAt = time.Now()

	token, err := a.GetToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token != "Bearer cached-token" {
		t.Fatalf("expected cached token")
	}
}

func TestIamApiKey_CacheExpired_TriggersRefreshPath(t *testing.T) {
	a, err := NewApiKey("dummy", false, false)
	if err != nil {
		t.Fatal(err)
	}

	// expired token
	a.token = "Bearer old"
	a.generatedAt = time.Now().Add(-31 * time.Minute)

	// This will attempt HTTP call → should fail
	_, err = a.GetToken(context.Background())

	if err == nil {
		t.Fatalf("expected error due to real HTTP call (no mocking)")
	}
}
