package services

import (
	"errors"

	"github.com/velostack/velostack-go/app/models"
	"github.com/velostack/velostack-go/app/repositories"
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

	// Verify old password
	if !checkPassword(user.Password, oldPassword) {
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
