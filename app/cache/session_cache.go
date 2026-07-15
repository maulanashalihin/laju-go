package cache

import (
	"encoding/json"
	"time"

	"github.com/nutsdb/nutsdb"
)

// CachedSessionData holds the session fields we cache.
// Survives restarts via NutsDB persistence.
type CachedSessionData struct {
	UserID     int64     `json:"uid"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	CSRFToken  string    `json:"csrf,omitempty"`
	CSRFExpiry int64     `json:"csrf_exp,omitempty"`
	IP         string    `json:"ip,omitempty"`
	UserAgent  string    `json:"ua,omitempty"`
	ExpiresAt  time.Time `json:"exp"`
}

// SessionCache provides NutsDB-backed session caching with TTL.
// Thread-safe via NutsDB transaction isolation.
type SessionCache struct {
	db  *nutsdb.DB
	ttl time.Duration
}

// NewSessionCache creates a session cache backed by NutsDB.
func NewSessionCache(db *nutsdb.DB, ttl time.Duration) *SessionCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &SessionCache{db: db, ttl: ttl}
}

// Get retrieves a cached session. Returns nil + false if not found or expired.
func (c *SessionCache) Get(sessionID string) (*CachedSessionData, bool) {
	var data CachedSessionData

	err := c.db.View(func(tx *nutsdb.Tx) error {
		val, err := tx.Get("sessions", []byte(sessionID))
		if err != nil {
			return err
		}
		return json.Unmarshal(val, &data)
	})
	if err != nil {
		return nil, false
	}

	// Check session expiry (from source of truth)
	if time.Now().After(data.ExpiresAt) {
		c.Invalidate(sessionID)
		return nil, false
	}

	return &data, true
}

// Set stores a session in NutsDB cache with TTL.
func (c *SessionCache) Set(sessionID string, data CachedSessionData) {
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Use session remaining TTL plus cache buffer as NutsDB TTL
	maxTTL := time.Until(data.ExpiresAt) + c.ttl
	if maxTTL < c.ttl {
		maxTTL = c.ttl
	}

	c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("sessions", []byte(sessionID), raw, ttlToUint32(maxTTL))
	})
}

// Invalidate removes a session from cache.
func (c *SessionCache) Invalidate(sessionID string) {
	c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete("sessions", []byte(sessionID))
	})
}

// Clear removes all cached sessions.
func (c *SessionCache) Clear() {
	c.db.Update(func(tx *nutsdb.Tx) error {
		keys, err := tx.GetKeys("sessions")
		if err != nil {
			return nil
		}
		for _, key := range keys {
			tx.Delete("sessions", key)
		}
		return nil
	})
}
