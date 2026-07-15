package cache

import (
	"os"
	"testing"
	"time"

	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestUserCache(t *testing.T, ttl time.Duration) (*UserCache, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "cache-test-*")
	require.NoError(t, err)

	ndb, err := Open(dir)
	require.NoError(t, err)

	c := NewUserCache(ndb.DB, ttl)
	return c, func() {
		ndb.Close()
		os.RemoveAll(dir)
	}
}

func TestUserCacheGetSet(t *testing.T) {
	c, cleanup := newTestUserCache(t, 5*time.Minute)
	defer cleanup()

	user := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	c.Set(user)

	got := c.Get(1)
	require.NotNil(t, got)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
}

func TestUserCacheGetMiss(t *testing.T) {
	c, cleanup := newTestUserCache(t, 5*time.Minute)
	defer cleanup()

	assert.Nil(t, c.Get(999))
}

func TestUserCacheInvalidate(t *testing.T) {
	c, cleanup := newTestUserCache(t, 5*time.Minute)
	defer cleanup()

	user := &models.User{ID: 1, Name: "Bob"}
	c.Set(user)
	assert.NotNil(t, c.Get(1))

	c.Invalidate(1)
	assert.Nil(t, c.Get(1))
}

func TestUserCacheExpiry(t *testing.T) {
	c, cleanup := newTestUserCache(t, 50*time.Millisecond)
	defer cleanup()

	user := &models.User{ID: 1, Name: "Carol"}
	c.Set(user)
	assert.NotNil(t, c.Get(1))

	time.Sleep(60 * time.Millisecond)
	assert.Nil(t, c.Get(1))
}

func TestUserCacheClear(t *testing.T) {
	c, cleanup := newTestUserCache(t, 5*time.Minute)
	defer cleanup()

	c.Set(&models.User{ID: 1, Name: "One"})
	c.Set(&models.User{ID: 2, Name: "Two"})
	assert.Equal(t, 2, c.Size())

	c.Clear()
	assert.Equal(t, 0, c.Size())
}

func TestUserCacheSize(t *testing.T) {
	c, cleanup := newTestUserCache(t, 5*time.Minute)
	defer cleanup()

	assert.Equal(t, 0, c.Size())
	c.Set(&models.User{ID: 1, Name: "A"})
	assert.Equal(t, 1, c.Size())
	c.Set(&models.User{ID: 2, Name: "B"})
	assert.Equal(t, 2, c.Size())
}

func TestUserCacheOverwrite(t *testing.T) {
	c, cleanup := newTestUserCache(t, 5*time.Minute)
	defer cleanup()

	c.Set(&models.User{ID: 1, Name: "Original"})
	c.Set(&models.User{ID: 1, Name: "Updated"})

	got := c.Get(1)
	require.NotNil(t, got)
	assert.Equal(t, "Updated", got.Name)
}
