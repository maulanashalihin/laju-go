package services

import (
	"errors"

	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/maulanashalihin/laju-go/app/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves a user's profile
func (s *UserService) GetProfile(userID int64) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetProfileByEmail retrieves a user's profile by email
func (s *UserService) GetProfileByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID int64, hashedPassword string) error {
	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

// UpdateAvatar updates a user's avatar URL
func (s *UserService) UpdateAvatar(userID int64, avatarURL string) error {
	return s.userRepo.UpdateAvatar(userID, avatarURL)
}

// UpdateProfile updates a user's profile
func (s *UserService) UpdateProfile(userID int64, req models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
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

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
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
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

// DeleteAccount deletes a user's account
func (s *UserService) DeleteAccount(userID int64) error {
	return s.userRepo.Delete(userID)
}

// IsAdmin checks if a user is an admin
func (s *UserService) IsAdmin(userID int64) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == models.RoleAdmin, nil
}
