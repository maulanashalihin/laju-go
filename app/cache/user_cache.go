package cache

import (
	"sync"
	"time"

	"github.com/maulanashalihin/laju-go/app/models"
)

type cacheEntry struct {
	user      *models.User
	expiresAt time.Time
}

// UserCache provides TTL-based in-memory caching for user profiles.
// Thread-safe via sync.RWMutex.
type UserCache struct {
	mu   sync.RWMutex
	data map[int64]cacheEntry
	ttl  time.Duration
}

// NewUserCache creates a user profile cache with the given TTL.
// A background goroutine periodically purges expired entries.
func NewUserCache(ttl time.Duration) *UserCache {
	c := &UserCache{
		data: make(map[int64]cacheEntry),
		ttl:  ttl,
	}
	go c.cleanup()
	return c
}

// Get retrieves a user from cache. Returns nil if not found or expired.
func (c *UserCache) Get(userID int64) *models.User {
	c.mu.RLock()
	entry, ok := c.data[userID]
	c.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		// Expired: clean up
		if ok {
			c.mu.Lock()
			delete(c.data, userID)
			c.mu.Unlock()
		}
		return nil
	}

	return entry.user
}

// Set stores a user in cache with the configured TTL.
func (c *UserCache) Set(user *models.User) {
	c.mu.Lock()
	c.data[user.ID] = cacheEntry{
		user:      user,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

// Invalidate removes a user from cache (call after updates).
func (c *UserCache) Invalidate(userID int64) {
	c.mu.Lock()
	delete(c.data, userID)
	c.mu.Unlock()
}

// Clear removes all cached entries.
func (c *UserCache) Clear() {
	c.mu.Lock()
	c.data = make(map[int64]cacheEntry)
	c.mu.Unlock()
}

// cleanup runs in a goroutine, evicting expired entries every 5 minutes.
func (c *UserCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for id, entry := range c.data {
			if now.After(entry.expiresAt) {
				delete(c.data, id)
			}
		}
		c.mu.Unlock()
	}
}

// Size returns the number of non-expired cached entries (for debugging).
func (c *UserCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	for _, entry := range c.data {
		if time.Now().Before(entry.expiresAt) {
			count++
		}
	}
	return count
}
