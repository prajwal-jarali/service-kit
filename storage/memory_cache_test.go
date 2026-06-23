package storage

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInMemoryCache_SetGet(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	err := c.Set(ctx, "k1", []byte("v1"), 0)
	if err != nil {
		t.Fatal(err)
	}

	val, ok, err := c.Get(ctx, "k1")
	if err != nil || !ok {
		t.Fatalf("expected key present")
	}

	if string(val) != "v1" {
		t.Fatalf("unexpected value")
	}
}

func TestInMemoryCache_Overwrite(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	c.Set(ctx, "k1", []byte("v1"), 0)
	c.Set(ctx, "k1", []byte("v2"), 0)

	val, ok, _ := c.Get(ctx, "k1")
	if !ok || string(val) != "v2" {
		t.Fatalf("expected overwritten value")
	}
}

func TestInMemoryCache_TTL_LazyExpiration(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	c.Set(ctx, "k1", []byte("v1"), 20*time.Millisecond)

	time.Sleep(50 * time.Millisecond)

	_, ok, _ := c.Get(ctx, "k1")
	if ok {
		t.Fatalf("expected key to expire")
	}
}

func TestInMemoryCache_Delete(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	c.Set(ctx, "k1", []byte("v1"), 0)
	c.Delete(ctx, "k1")

	_, ok, _ := c.Get(ctx, "k1")
	if ok {
		t.Fatalf("expected key deleted")
	}
}

func TestInMemoryCache_Clear(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	c.Set(ctx, "k1", []byte("v1"), 0)
	c.Set(ctx, "k2", []byte("v2"), 0)

	c.Clear(ctx)

	_, ok1, _ := c.Get(ctx, "k1")
	_, ok2, _ := c.Get(ctx, "k2")

	if ok1 || ok2 {
		t.Fatalf("expected cache to be cleared")
	}
}

func TestInMemoryCache_Exists(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	c.Set(ctx, "k1", []byte("v1"), 0)

	ok, err := c.Exists(ctx, "k1")
	if err != nil || !ok {
		t.Fatalf("expected exists true")
	}
}

func TestInMemoryCache_Concurrency(t *testing.T) {
	c := NewInMemoryCache(0)
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			key := "k"

			c.Set(ctx, key, []byte("v"), 0)
			c.Get(ctx, key)
			c.Exists(ctx, key)
		}(i)
	}

	wg.Wait()

	// final sanity check
	_, ok, _ := c.Get(ctx, "k")
	if !ok {
		t.Fatalf("expected key to exist after concurrent ops")
	}
}
