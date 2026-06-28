package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/cache"
	"github.com/maulanashalihin/laju-go/app/queries"
)

type Store struct {
	querier      *queries.Querier
	sessionCache *cache.SessionCache
	sessionName  string
	sessionTTL   time.Duration
}

type Session struct {
	id        string
	userID    int64
	values    map[string]interface{}
	c         *fiber.Ctx
	store     *Store
	dirty     bool
	expiresAt time.Time
}

// SessionData represents the data stored in session
type SessionData struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	CSRFToken   string `json:"csrf_token,omitempty"`
	CSRFExpiry  int64  `json:"csrf_expiry,omitempty"`
}

// New creates a new session store with database backend and optional in-memory cache.
// Pass nil for sessionCache to disable caching. sessionTTL is the session lifetime.
func New(querier *queries.Querier, sessionCache *cache.SessionCache, sessionTTL time.Duration) *Store {
	if sessionTTL <= 0 {
		sessionTTL = 24 * time.Hour
	}
	return &Store{
		querier:      querier,
		sessionCache: sessionCache,
		sessionName:  "session_id",
		sessionTTL:   sessionTTL,
	}
}

// Get retrieves a session
func (s *Store) Get(c *fiber.Ctx) (*Session, error) {
	// Get session from locals first (if already loaded)
	if sess := c.Locals("session"); sess != nil {
		slog.Info("session retrieved from locals")
		return sess.(*Session), nil
	}

	session := &Session{
		id:        "",
		userID:    0,
		values:    make(map[string]interface{}),
		c:         c,
		store:     s,
		dirty:     false,
		expiresAt: time.Now().Add(s.sessionTTL),
	}

	// Try to get existing session from cookie
	cookieValue := c.Cookies(s.sessionName)
	if cookieValue != "" {
		// Try in-memory cache first (avoids DB lookup on every request)
		if s.sessionCache != nil {
			if cached, ok := s.sessionCache.Get(cookieValue); ok {
				// Check if the DB session has expired (cache TTL may outlive session TTL)
				if cached.ExpiresAt.Before(time.Now()) {
					s.sessionCache.Invalidate(cookieValue)
				} else {
					session.id = cookieValue
					session.userID = cached.UserID
					session.expiresAt = cached.ExpiresAt
					session.values["user_id"] = cached.UserID
					session.values["email"] = cached.Email
					session.values["role"] = cached.Role
					if cached.CSRFToken != "" {
						session.values["csrf_token"] = cached.CSRFToken
						session.values["csrf_expiry"] = cached.CSRFExpiry
					}
					c.Locals("session", session)
					return session, nil
				}
			}
		}

		// Cache miss: find session in database
		dbSession, err := s.querier.GetSessionByID(context.Background(), cookieValue)
		if err == nil {
			// Check if session is expired
			if dbSession.ExpiresAt.Before(time.Now()) {
				s.querier.DeleteSession(context.Background(), cookieValue)
				if s.sessionCache != nil {
					s.sessionCache.Invalidate(cookieValue)
				}
			} else {
				session.id = dbSession.ID
				session.userID = dbSession.UserID
				session.expiresAt = dbSession.ExpiresAt

				// Decode session data
				var data SessionData
				if err := json.Unmarshal([]byte(dbSession.Data), &data); err == nil {
					session.values["user_id"] = data.UserID
					session.values["email"] = data.Email
					session.values["role"] = data.Role
					if data.CSRFToken != "" {
						session.values["csrf_token"] = data.CSRFToken
					}
					if data.CSRFExpiry != 0 {
						session.values["csrf_expiry"] = data.CSRFExpiry
					}

					// Store in cache for subsequent requests
					if s.sessionCache != nil {
						s.sessionCache.Set(cookieValue, cache.CachedSessionData{
							UserID:     data.UserID,
							Email:      data.Email,
							Role:       data.Role,
							CSRFToken:  data.CSRFToken,
							CSRFExpiry: data.CSRFExpiry,
							ExpiresAt:  dbSession.ExpiresAt,
						})
					}
				}
			}
		} else {
			// Session not found or expired in DB — invalidate stale cache entry
			if s.sessionCache != nil {
				s.sessionCache.Invalidate(cookieValue)
			}
		}
	}

	c.Locals("session", session)
	return session, nil
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Set sets a value in the session
func (s *Session) Set(key string, value interface{}) {
	s.values[key] = value
	s.dirty = true
}

// Get gets a value from the session
func (s *Session) Get(key string) interface{} {
	return s.values[key]
}

// Delete removes a value from the session
func (s *Session) Delete(key string) {
	delete(s.values, key)
	s.dirty = true
}

// Save saves the session to database
func (s *Session) Save() error {
	sessionData := SessionData{}

	if userID, ok := s.values["user_id"].(int64); ok {
		sessionData.UserID = userID
	} else if userID, ok := s.values["user_id"].(int); ok {
		sessionData.UserID = int64(userID)
	} else if userID, ok := s.values["user_id"].(float64); ok {
		sessionData.UserID = int64(userID)
	}

	if email, ok := s.values["email"].(string); ok {
		sessionData.Email = email
	}

	if role, ok := s.values["role"].(string); ok {
		sessionData.Role = role
	}

	if csrfToken, ok := s.values["csrf_token"].(string); ok {
		sessionData.CSRFToken = csrfToken
	}

	if csrfExpiry, ok := s.values["csrf_expiry"].(int64); ok {
		sessionData.CSRFExpiry = csrfExpiry
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		slog.Error("session marshal error", "error", err)
		return err
	}

	// Refresh expiry on every save (sliding expiration)
	s.expiresAt = time.Now().Add(s.store.sessionTTL)

	if s.id == "" {
		// Create new session
		sessionID, err := generateSessionID()
		if err != nil {
			slog.Error("session generate id error", "error", err)
			return err
		}
		s.id = sessionID

		dbSession := &queries.Session{
			ID:        s.id,
			UserID:    sessionData.UserID,
			Data:      string(jsonData),
			ExpiresAt: s.expiresAt,
		}

		if err := s.store.querier.CreateSession(context.Background(), dbSession); err != nil {
			slog.Error("session create error", "error", err)
			return err
		}

		// Seed cache after creation
		if s.store.sessionCache != nil {
			s.store.sessionCache.Set(s.id, cache.CachedSessionData{
				UserID:     sessionData.UserID,
				Email:      sessionData.Email,
				Role:       sessionData.Role,
				CSRFToken:  sessionData.CSRFToken,
				CSRFExpiry: sessionData.CSRFExpiry,
				ExpiresAt:  s.expiresAt,
			})
		}
	} else {
		// Update existing session
		dbSession := &queries.Session{
			ID:        s.id,
			UserID:    sessionData.UserID,
			Data:      string(jsonData),
			ExpiresAt: s.expiresAt,
		}

		if err := s.store.querier.UpdateSession(context.Background(), dbSession); err != nil {
			slog.Error("session update error", "error", err)
			return err
		}

		// Refresh cache after update
		if s.store.sessionCache != nil {
			s.store.sessionCache.Set(s.id, cache.CachedSessionData{
				UserID:     sessionData.UserID,
				Email:      sessionData.Email,
				Role:       sessionData.Role,
				CSRFToken:  sessionData.CSRFToken,
				CSRFExpiry: sessionData.CSRFExpiry,
				ExpiresAt:  s.expiresAt,
			})
		}
	}

	// Set cookie with session ID
	s.c.Cookie(&fiber.Cookie{
		Name:     s.store.sessionName,
		Value:    s.id,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		MaxAge:   int(s.expiresAt.Sub(time.Now()).Seconds()),
	})
	slog.Info("session cookie set", "name", s.store.sessionName, "value", s.id)

	return nil
}

// Destroy destroys the session
func (s *Session) Destroy() error {
	if s.id != "" {
		s.store.querier.DeleteSession(context.Background(), s.id)
		if s.store.sessionCache != nil {
			s.store.sessionCache.Invalidate(s.id)
		}
	}

	s.values = make(map[string]interface{})
	s.c.ClearCookie(s.store.sessionName)
	return nil
}

// Regenerate generates a new session ID
func (s *Session) Regenerate() error {
	if s.id == "" {
		return nil // Nothing to regenerate
	}

	newID, err := generateSessionID()
	if err != nil {
		return err
	}

	// Re-encode data
	sessionData := SessionData{}

	if userID, ok := s.values["user_id"].(int64); ok {
		sessionData.UserID = userID
	} else if userID, ok := s.values["user_id"].(int); ok {
		sessionData.UserID = int64(userID)
	} else if userID, ok := s.values["user_id"].(float64); ok {
		sessionData.UserID = int64(userID)
	}

	if email, ok := s.values["email"].(string); ok {
		sessionData.Email = email
	}

	if role, ok := s.values["role"].(string); ok {
		sessionData.Role = role
	}

	if csrfToken, ok := s.values["csrf_token"].(string); ok {
		sessionData.CSRFToken = csrfToken
	}

	if csrfExpiry, ok := s.values["csrf_expiry"].(int64); ok {
		sessionData.CSRFExpiry = csrfExpiry
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	// Create new session
	dbSession := &queries.Session{
		ID:        newID,
		UserID:    sessionData.UserID,
		Data:      string(jsonData),
		ExpiresAt: s.expiresAt,
	}

	if err := s.store.querier.CreateSession(context.Background(), dbSession); err != nil {
		return err
	}

	s.store.querier.DeleteSession(context.Background(), s.id)
	if s.store.sessionCache != nil {
		s.store.sessionCache.Invalidate(s.id)
		s.store.sessionCache.Set(newID, cache.CachedSessionData{
			UserID:     sessionData.UserID,
			Email:      sessionData.Email,
			Role:       sessionData.Role,
			CSRFToken:  sessionData.CSRFToken,
			CSRFExpiry: sessionData.CSRFExpiry,
			ExpiresAt:  s.expiresAt,
		})
	}

	s.id = newID

	// Update cookie
	s.c.Cookie(&fiber.Cookie{
		Name:     s.store.sessionName,
		Value:    s.id,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		MaxAge:   int(s.expiresAt.Sub(time.Now()).Seconds()),
	})

	return nil
}

// Flash sets a flash message cookie (short-lived, one-time use)
// The flash message will be available on the next request and then cleared
func (s *Store) Flash(c *fiber.Ctx, key string, value string) {
	// Set flash cookie with short expiry (5 minutes)
	c.Cookie(&fiber.Cookie{
		Name:     "flash_" + key,
		Value:    value,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		MaxAge:   300, // 5 minutes
	})
}

// GetFlash retrieves and clears a flash message cookie
func (s *Store) GetFlash(c *fiber.Ctx, key string) string {
	cookieName := "flash_" + key
	value := c.Cookies(cookieName)
	
	if value != "" {
		// Clear the flash cookie after reading (one-time use)
		c.ClearCookie(cookieName)
	}
	
	return value
}
