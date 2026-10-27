package alerter

import (
	"fmt"
	"sync"
	"testing"
)

func TestCache_IsNew(t *testing.T) {
	cache := NewCache(10)

	// First time seeing ID 1 -> should be new (true)
	if !cache.IsNew("sneakers", 1) {
		t.Errorf("expected ID 1 to be new, got false")
	}

	// Second time seeing ID 1 -> should NOT be new (false)
	if cache.IsNew("sneakers", 1) {
		t.Errorf("expected ID 1 to not be new on second check, got true")
	}

	// Different query with same ID -> should be considered new for that query
	if !cache.IsNew("hoodies", 1) {
		t.Errorf("expected ID 1 in different query to be new, got false")
	}
}

func TestCache_Concurrency(t *testing.T) {
	cache := NewCache(1000)
	var wg sync.WaitGroup
	workers := 50
	itemsPerWorker := 100

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < itemsPerWorker; i++ {
				id := workerID*itemsPerWorker + i
				query := fmt.Sprintf("query_%d", workerID%5)
				cache.IsNew(query, id)
			}
		}(w)
	}

	wg.Wait()
}

func TestCache_Eviction(t *testing.T) {
	maxLen := 5
	cache := NewCache(maxLen)

	// Fill cache to max
	for i := 1; i <= maxLen; i++ {
		if !cache.IsNew("test", i) {
			t.Fatalf("expected item %d to be new", i)
		}
	}

	// Insert item that causes reset/eviction
	if !cache.IsNew("test", maxLen+1) {
		t.Fatalf("expected item %d to be new", maxLen+1)
	}

	cache.mu.RLock()
	queryLen := len(cache.items["test"])
	cache.mu.RUnlock()

	if queryLen > maxLen {
		t.Errorf("cache size %d exceeded maxLen %d", queryLen, maxLen)
	}
}
