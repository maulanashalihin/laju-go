package session

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/securecookie"
)

type Store struct {
	secureCookie *securecookie.SecureCookie
	sessionName  string
}

type Session struct {
	values map[interface{}]interface{}
	c      *fiber.Ctx
	store  *Store
	dirty  bool
}

// New creates a new session store
func New(secret string) *Store {
	return &Store{
		secureCookie: securecookie.New([]byte(secret), nil),
		sessionName:  "session_id",
	}
}

// Get retrieves a session
func (s *Store) Get(c *fiber.Ctx) (*Session, error) {
	// Get session from locals first (if already loaded)
	if sess := c.Locals("session"); sess != nil {
		return sess.(*Session), nil
	}

	session := &Session{
		values: make(map[interface{}]interface{}),
		c:      c,
		store:  s,
		dirty:  false,
	}

	// Try to get existing session from cookie
	cookieValue := c.Cookies(s.sessionName)
	if cookieValue != "" {
		var data map[interface{}]interface{}
		if err := s.secureCookie.Decode(s.sessionName, cookieValue, &data); err == nil {
			session.values = data
		}
	}

	c.Locals("session", session)
	return session, nil
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

// Save saves the session
func (s *Session) Save() error {
	if !s.dirty {
		return nil
	}

	encoded, err := s.store.secureCookie.Encode(s.store.sessionName, s.values)
	if err != nil {
		return err
	}

	s.c.Cookie(&fiber.Cookie{
		Name:     s.store.sessionName,
		Value:    encoded,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return nil
}

// Destroy destroys the session
func (s *Session) Destroy() error {
	s.values = make(map[interface{}]interface{})
	s.c.ClearCookie(s.store.sessionName)
	return nil
}

// Helper function for base64 encoding/decoding (alternative to securecookie)
func encodeBase64(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func decodeBase64(encoded string, data interface{}) error {
	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, data)
}
