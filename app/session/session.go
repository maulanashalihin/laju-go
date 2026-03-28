package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/maulanashalihin/laju-go/app/repositories"
)

type Store struct {
	repo        *repositories.SessionRepository
	sessionName string
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
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// New creates a new session store with database backend
func New(repo *repositories.SessionRepository) *Store {
	return &Store{
		repo:        repo,
		sessionName: "session_id",
	}
}

// Get retrieves a session
func (s *Store) Get(c *fiber.Ctx) (*Session, error) {
	// Get session from locals first (if already loaded)
	if sess := c.Locals("session"); sess != nil {
		log.Printf("[Session] Retrieved from locals\n")
		return sess.(*Session), nil
	}

	session := &Session{
		id:        "",
		userID:    0,
		values:    make(map[string]interface{}),
		c:         c,
		store:     s,
		dirty:     false,
		expiresAt: time.Now().Add(24 * time.Hour), // Default 24 hours
	}

	// Try to get existing session from cookie
	cookieValue := c.Cookies(s.sessionName)
	log.Printf("[Session] Cookie value: '%s'\n", cookieValue)
	
	if cookieValue != "" {
		// Find session in database
		dbSession, err := s.repo.GetByID(cookieValue)
		if err == nil {
			// Session found in database
			session.id = dbSession.ID
			session.userID = dbSession.UserID
			session.expiresAt = dbSession.ExpiresAt

			// Decode session data
			var data SessionData
			if err := json.Unmarshal([]byte(dbSession.Data), &data); err == nil {
				session.values["user_id"] = data.UserID
				session.values["email"] = data.Email
				session.values["role"] = data.Role
				log.Printf("[Session] Loaded from DB: id=%s, user_id=%d\n", session.id, data.UserID)
			} else {
				log.Printf("[Session] Decode error: %v\n", err)
			}
		} else {
			log.Printf("[Session] DB lookup error: %v\n", err)
		}
		// If session not found or expired, start fresh
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
	// Encode session data
	sessionData := SessionData{
		UserID: 0,
		Email:  "",
		Role:   "",
	}

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

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		log.Printf("[Session] Marshal error: %v\n", err)
		return err
	}

	if s.id == "" {
		// Create new session
		sessionID, err := generateSessionID()
		if err != nil {
			log.Printf("[Session] Generate ID error: %v\n", err)
			return err
		}
		s.id = sessionID

		dbSession := &models.Session{
			ID:        s.id,
			UserID:    sessionData.UserID,
			Data:      string(jsonData),
			ExpiresAt: s.expiresAt,
		}

		if err := s.store.repo.Create(dbSession); err != nil {
			log.Printf("[Session] Create error: %v\n", err)
			return err
		}
		log.Printf("[Session] Created new session: id=%s, user_id=%d\n", s.id, sessionData.UserID)
	} else {
		// Update existing session
		dbSession := &models.Session{
			ID:        s.id,
			UserID:    sessionData.UserID,
			Data:      string(jsonData),
			ExpiresAt: s.expiresAt,
		}

		if err := s.store.repo.Update(dbSession); err != nil {
			log.Printf("[Session] Update error: %v\n", err)
			return err
		}
		log.Printf("[Session] Updated session: id=%s, user_id=%d\n", s.id, sessionData.UserID)
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
	log.Printf("[Session] Cookie set: name=%s, value=%s\n", s.store.sessionName, s.id)

	return nil
}

// Destroy destroys the session
func (s *Session) Destroy() error {
	if s.id != "" {
		// Delete from database
		s.store.repo.Delete(s.id)
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

	// Update session ID in database
	dbSession := &models.Session{
		ID:        newID,
		UserID:    s.userID,
		Data:      "", // Will be re-encoded
		ExpiresAt: s.expiresAt,
	}

	// Re-encode data
	sessionData := SessionData{
		UserID: 0,
		Email:  "",
		Role:   "",
	}

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

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	dbSession.Data = string(jsonData)

	// Create new session
	if err := s.store.repo.Create(dbSession); err != nil {
		return err
	}

	// Delete old session
	s.store.repo.Delete(s.id)

	// Update local ID
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
