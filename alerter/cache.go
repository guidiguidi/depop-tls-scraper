package alerter

import (
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	items  map[string]map[int]bool
	maxLen int
}

// NewCache initializes an in-memory cache
func NewCache(maxLen int) *Cache {
	return &Cache{
		items:  make(map[string]map[int]bool),
		maxLen: maxLen,
	}
}

// IsNew checks if an item is new and adds it to the cache if it is.
// Returns true if the item was NOT in the cache (i.e. it's new).
func (c *Cache) IsNew(query string, id int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[query]; !ok {
		c.items[query] = make(map[int]bool)
	}

	queryCache := c.items[query]

	if queryCache[id] {
		return false // Not new
	}

	// Add to cache
	queryCache[id] = true

	// Extremely simple size limiting logic to prevent memory leaks.
	// Since items are usually sorted by newlyListed, we just reset the cache
	// if it gets too large. A real LRU cache would be better for production,
	// but this handles the basic requirement of bounded memory.
	if len(queryCache) > c.maxLen {
		// Just clear it and keep the current one.
		// A better approach in a real app is to delete oldest entries.
		c.items[query] = make(map[int]bool)
		c.items[query][id] = true
	}

	return true // It's new
}

// IsColdStart checks if a query cache is completely empty
func (c *Cache) IsColdStart(query string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if qc, ok := c.items[query]; ok {
		return len(qc) == 0
	}
	return true
}
