package cache

import (
	"sync"
	"time"
)

// CachedSessionData holds the session fields we cache in memory.
// This avoids a database lookup on every authenticated request.
type CachedSessionData struct {
	UserID     int64
	Email      string
	Role       string
	CSRFToken  string
	CSRFExpiry int64
}

type sessionCacheEntry struct {
	data      CachedSessionData
	expiresAt time.Time
}

// SessionCache provides TTL-based in-memory caching for sessions.
// Thread-safe via sync.RWMutex.
type SessionCache struct {
	mu   sync.RWMutex
	data map[string]sessionCacheEntry
	ttl  time.Duration
}

// NewSessionCache creates a session cache with the given TTL.
// A background goroutine periodically purges expired entries.
func NewSessionCache(ttl time.Duration) *SessionCache {
	c := &SessionCache{
		data: make(map[string]sessionCacheEntry),
		ttl:  ttl,
	}
	go c.cleanup()
	return c
}

// Get retrieves a cached session. Returns nil + false if not found or expired.
func (c *SessionCache) Get(sessionID string) (*CachedSessionData, bool) {
	c.mu.RLock()
	entry, ok := c.data[sessionID]
	c.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		if ok {
			c.mu.Lock()
			delete(c.data, sessionID)
			c.mu.Unlock()
		}
		return nil, false
	}

	return &entry.data, true
}

// Set stores a session in cache with the configured TTL.
func (c *SessionCache) Set(sessionID string, data CachedSessionData) {
	c.mu.Lock()
	c.data[sessionID] = sessionCacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

// Invalidate removes a session from cache (call after update/destroy).
func (c *SessionCache) Invalidate(sessionID string) {
	c.mu.Lock()
	delete(c.data, sessionID)
	c.mu.Unlock()
}

// Clear removes all cached sessions.
func (c *SessionCache) Clear() {
	c.mu.Lock()
	c.data = make(map[string]sessionCacheEntry)
	c.mu.Unlock()
}

// cleanup runs in a goroutine, evicting expired entries every 5 minutes.
func (c *SessionCache) cleanup() {
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
