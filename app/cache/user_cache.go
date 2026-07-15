package cache

import (
	"encoding/json"
	"time"

	"github.com/nutsdb/nutsdb"
	"github.com/maulanashalihin/laju-go/app/models"
)

// userCacheEntry wraps a user with application-level TTL.
type userCacheEntry struct {
	User      *models.User `json:"u"`
	ExpiresAt time.Time    `json:"exp"`
}

// UserCache provides NutsDB-backed user profile caching with TTL.
// Thread-safe via NutsDB transaction isolation.
type UserCache struct {
	db  *nutsdb.DB
	ttl time.Duration
}

// NewUserCache creates a user profile cache backed by NutsDB.
func NewUserCache(db *nutsdb.DB, ttl time.Duration) *UserCache {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &UserCache{db: db, ttl: ttl}
}

// Get retrieves a user from cache. Returns nil if not found or expired.
func (c *UserCache) Get(userID int64) *models.User {
	if c.db == nil {
		return nil
	}

	var entry userCacheEntry

	err := c.db.View(func(tx *nutsdb.Tx) error {
		val, err := tx.Get("users", int64Key(userID))
		if err != nil {
			return err
		}
		return json.Unmarshal(val, &entry)
	})
	if err != nil {
		return nil
	}

	// Application-level TTL check (supports sub-second precision)
	if time.Now().After(entry.ExpiresAt) {
		c.Invalidate(userID)
		return nil
	}

	return entry.User
}

// Set stores a user in cache with the configured TTL.
func (c *UserCache) Set(user *models.User) {
	if c.db == nil || user == nil {
		return
	}

	entry := userCacheEntry{
		User:      user,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	raw, err := json.Marshal(entry)
	if err != nil {
		return
	}

	// NutsDB native TTL (second granularity) as safety net
	c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("users", int64Key(user.ID), raw, ttlToUint32(c.ttl))
	})
}

// Invalidate removes a user from cache (call after updates).
func (c *UserCache) Invalidate(userID int64) {
	if c.db == nil {
		return
	}
	c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete("users", int64Key(userID))
	})
}

// Clear removes all cached entries.
func (c *UserCache) Clear() {
	if c.db == nil {
		return
	}
	c.db.Update(func(tx *nutsdb.Tx) error {
		keys, err := tx.GetKeys("users")
		if err != nil {
			return nil
		}
		for _, key := range keys {
			tx.Delete("users", key)
		}
		return nil
	})
}

// Size returns the approximate number of cached entries.
func (c *UserCache) Size() int {
	if c.db == nil {
		return 0
	}
	count := 0
	c.db.View(func(tx *nutsdb.Tx) error {
		keys, err := tx.GetKeys("users")
		if err != nil {
			return nil
		}
		count = len(keys)
		return nil
	})
	return count
}
