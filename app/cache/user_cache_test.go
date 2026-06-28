package cache

import (
	"testing"
	"time"

	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCacheGetSet(t *testing.T) {
	c := NewUserCache(5 * time.Minute)

	user := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	c.Set(user)

	got := c.Get(1)
	require.NotNil(t, got)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "alice@example.com", got.Email)
}

func TestUserCacheGetMiss(t *testing.T) {
	c := NewUserCache(5 * time.Minute)
	assert.Nil(t, c.Get(999))
}

func TestUserCacheInvalidate(t *testing.T) {
	c := NewUserCache(5 * time.Minute)

	user := &models.User{ID: 1, Name: "Bob"}
	c.Set(user)
	assert.NotNil(t, c.Get(1))

	c.Invalidate(1)
	assert.Nil(t, c.Get(1))
}

func TestUserCacheExpiry(t *testing.T) {
	c := NewUserCache(50 * time.Millisecond)

	user := &models.User{ID: 1, Name: "Carol"}
	c.Set(user)
	assert.NotNil(t, c.Get(1))

	time.Sleep(60 * time.Millisecond)
	assert.Nil(t, c.Get(1))
}

func TestUserCacheClear(t *testing.T) {
	c := NewUserCache(5 * time.Minute)

	c.Set(&models.User{ID: 1, Name: "One"})
	c.Set(&models.User{ID: 2, Name: "Two"})
	assert.Equal(t, 2, c.Size())

	c.Clear()
	assert.Equal(t, 0, c.Size())
}

func TestUserCacheSize(t *testing.T) {
	c := NewUserCache(5 * time.Minute)

	assert.Equal(t, 0, c.Size())
	c.Set(&models.User{ID: 1, Name: "A"})
	assert.Equal(t, 1, c.Size())
	c.Set(&models.User{ID: 2, Name: "B"})
	assert.Equal(t, 2, c.Size())
}

func TestUserCacheOverwrite(t *testing.T) {
	c := NewUserCache(5 * time.Minute)

	c.Set(&models.User{ID: 1, Name: "Original"})
	c.Set(&models.User{ID: 1, Name: "Updated"})

	got := c.Get(1)
	require.NotNil(t, got)
	assert.Equal(t, "Updated", got.Name)
}
