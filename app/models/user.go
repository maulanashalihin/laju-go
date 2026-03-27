package models

import "time"

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	Password       string    `json:"-"` // Hashed password, never return in JSON
	Role           UserRole  `json:"role"`
	GoogleID       string    `json:"-"` // OAuth provider ID
	EmailVerified  bool      `json:"email_verified"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}
