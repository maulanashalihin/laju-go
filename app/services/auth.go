package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/maulanashalihin/laju-go/app/models"
	"github.com/maulanashalihin/laju-go/app/repositories"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	userRepo      *repositories.UserRepository
	sessionSecret string
	oauthConfig   *oauth2.Config
}

type AuthServiceConfig struct {
	SessionSecret      string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func NewAuthService(userRepo *repositories.UserRepository, cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		sessionSecret: cfg.SessionSecret,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

// GetOAuthConfig returns the OAuth config for Google
func (s *AuthService) GetOAuthConfig() *oauth2.Config {
	return s.oauthConfig
}

// ProcessGoogleToken exchanges the OAuth code for a token and returns user info
func (s *AuthService) ProcessGoogleToken(ctx context.Context, code string) (*models.User, error) {
	// Exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user info from Google
	oauthClient := s.oauthConfig.Client(ctx, token)
	resp, err := oauthClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, ErrInvalidToken
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Verified bool   `json:"verified_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, ErrInvalidToken
	}

	// Check if user exists by Google ID
	user, err := s.userRepo.GetByGoogleID(googleUser.ID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, err
	}

	// Check if user exists by email
	user, err = s.userRepo.GetByEmail(googleUser.Email)
	if err == nil {
		// Link Google ID to existing account
		user.GoogleID = sql.NullString{String: googleUser.ID, Valid: true}
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// Create new user
	newUser := &models.User{
		Email: googleUser.Email,
		Name:  googleUser.Name,
		GoogleID: sql.NullString{
			String: googleUser.ID,
			Valid:  true,
		},
		Avatar:        googleUser.Picture,
		EmailVerified: googleUser.Verified,
		Role:          models.RoleUser,
	}

	if err := s.userRepo.CreateWithGoogleID(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Register creates a new user with email/password
func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, repositories.ErrUserAlreadyExists
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email: email,
		Name:  name,
		Password: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Role:          models.RoleUser,
		EmailVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user with email/password
func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password - user must have a password (not OAuth-only user)
	if !user.Password.Valid {
		return nil, ErrInvalidCredentials
	}

	if !checkPassword(user.Password.String, password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id int64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword checks if a password matches a hash
func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetOAuthURL returns the OAuth URL for Google login
func (s *AuthService) GetOAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ValidateState validates the OAuth state parameter
func (s *AuthService) ValidateState(state, expected string) bool {
	return state == expected
}
