package storage

import (
	"context"
	"testing"
	"time"
)

// compile-time check
func assertCache(c Cache) {}

func TestCache_InterfaceImplementation(t *testing.T) {
	c := NewInMemoryCache(0) // no cleanup goroutine
	assertCache(c)
}

func TestCache_BasicContract(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	// Set
	err := c.Set(ctx, "k1", []byte("v1"), 0)
	if err != nil {
		t.Fatal(err)
	}

	// Get
	val, ok, err := c.Get(ctx, "k1")
	if err != nil || !ok {
		t.Fatalf("expected key present")
	}
	if string(val) != "v1" {
		t.Fatalf("unexpected value")
	}

	// Exists
	exists, err := c.Exists(ctx, "k1")
	if err != nil || !exists {
		t.Fatalf("expected exists true")
	}

	// Delete
	err = c.Delete(ctx, "k1")
	if err != nil {
		t.Fatal(err)
	}

	_, ok, _ = c.Get(ctx, "k1")
	if ok {
		t.Fatalf("expected key deleted")
	}
}

func TestCache_TTLContract(t *testing.T) {
	c := NewInMemoryCache(0) // rely on lazy expiration
	ctx := context.Background()

	err := c.Set(ctx, "k1", []byte("v1"), 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	_, ok, _ := c.Get(ctx, "k1")
	if ok {
		t.Fatalf("expected key to expire")
	}
}
