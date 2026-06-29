package services

import (
	"context"
	"errors"

	"github.com/maulanashalihin/laju-go/app/cache"
	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/maulanashalihin/laju-go/app/queries"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	querier *queries.Querier
	cache   *cache.UserCache
}

func NewUserService(querier *queries.Querier, userCache *cache.UserCache) *UserService {
	return &UserService{
		querier: querier,
		cache:   userCache,
	}
}

// GetProfile retrieves a user's profile (with cache)
func (s *UserService) GetProfile(userID int64) (*models.UserResponse, error) {
	// Check cache first
	if user := s.cache.Get(userID); user != nil {
		response := user.ToResponse()
		return &response, nil
	}

	// Cache miss: query DB
	user, err := s.querier.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.cache.Set(user)

	response := user.ToResponse()
	return &response, nil
}

// GetProfileByEmail retrieves a user's profile by email
func (s *UserService) GetProfileByEmail(email string) (*models.User, error) {
	return s.querier.GetUserByEmail(context.Background(), email)
}

// UpdatePassword updates a user's password (invalidates cache)
func (s *UserService) UpdatePassword(userID int64, hashedPassword string) error {
	err := s.querier.UpdateUserPassword(context.Background(), userID, hashedPassword)
	if err == nil {
		s.cache.Invalidate(userID)
	}
	return err
}

// UpdateAvatar updates a user's avatar URL (invalidates cache)
func (s *UserService) UpdateAvatar(userID int64, avatarURL string) error {
	err := s.querier.UpdateUserAvatar(context.Background(), userID, avatarURL)
	if err == nil {
		s.cache.Invalidate(userID)
	}
	return err
}

// UpdateProfile updates a user's profile (invalidates cache)
func (s *UserService) UpdateProfile(userID int64, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.querier.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.querier.UpdateUser(context.Background(), user); err != nil {
		return nil, err
	}

	// Invalidate old cache entry
	s.cache.Invalidate(userID)

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword changes a user's password (invalidates cache)
func (s *UserService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.querier.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	// Verify old password - user must have a password
	if !user.Password.Valid {
		return errors.New("invalid current password")
	}

	if !checkPassword(user.Password.String, oldPassword) {
		return errors.New("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.querier.UpdateUserPassword(context.Background(), userID, string(hashedPassword))
	if err == nil {
		s.cache.Invalidate(userID)
	}
	return err
}

// DeleteAccount deletes a user's account (invalidates cache)
func (s *UserService) DeleteAccount(userID int64) error {
	err := s.querier.DeleteUser(context.Background(), userID)
	if err == nil {
		s.cache.Invalidate(userID)
	}
	return err
}

// IsAdmin checks if a user is an admin (with cache)
func (s *UserService) IsAdmin(userID int64) (bool, error) {
	if user := s.cache.Get(userID); user != nil {
		return user.Role == models.RoleAdmin, nil
	}

	user, err := s.querier.GetUserByID(context.Background(), userID)
	if err != nil {
		return false, err
	}

	s.cache.Set(user)
	return user.Role == models.RoleAdmin, nil
}
