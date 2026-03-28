package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/maulanashalihin/laju-go/app/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	db *sql.DB
	psql squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query, args, err := r.psql.
		Insert("users").
		Columns("email", "name", "password", "role", "created_at", "updated_at").
		Values(user.Email, user.Name, user.Password.String, user.Role, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

// CreateWithGoogleID creates a new user with Google OAuth
func (r *UserRepository) CreateWithGoogleID(user *models.User) error {
	query, args, err := r.psql.
		Insert("users").
		Columns("email", "name", "google_id", "avatar", "email_verified", "role", "created_at", "updated_at").
		Values(user.Email, user.Name, user.GoogleID, user.Avatar, user.EmailVerified, user.Role, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID finds a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query, args, err := r.psql.
		Select("id", "email", "name", "password", "avatar", "role", "google_id", "email_verified", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	var password sql.NullString
	var googleID sql.NullString
	err = r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Email, &user.Name, &password, &user.Avatar,
		&user.Role, &googleID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	user.Password = password
	user.GoogleID = googleID

	return user, nil
}

// GetByEmail finds a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query, args, err := r.psql.
		Select("id", "email", "name", "password", "avatar", "role", "google_id", "email_verified", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	var password sql.NullString
	var googleID sql.NullString
	err = r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Email, &user.Name, &password, &user.Avatar,
		&user.Role, &googleID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	user.Password = password
	user.GoogleID = googleID

	return user, nil
}

// GetByGoogleID finds a user by Google OAuth ID
func (r *UserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	query, args, err := r.psql.
		Select("id", "email", "name", "password", "avatar", "role", "google_id", "email_verified", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"google_id": googleID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	var password sql.NullString
	var googleIDNull sql.NullString
	err = r.db.QueryRow(query, args...).Scan(
		&user.ID, &user.Email, &user.Name, &password, &user.Avatar,
		&user.Role, &googleIDNull, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	user.Password = password
	user.GoogleID = googleIDNull

	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	query, args, err := r.psql.
		Update("users").
		Set("name", user.Name).
		Set("avatar", user.Avatar).
		Set("email_verified", user.EmailVerified).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id int64, hashedPassword string) error {
	query, args, err := r.psql.
		Update("users").
		Set("password", hashedPassword).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateAvatar updates a user's avatar URL
func (r *UserRepository) UpdateAvatar(id int64, avatarURL string) error {
	query, args, err := r.psql.
		Update("users").
		Set("avatar", avatarURL).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id int64) error {
	query, args, err := r.psql.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// SetAdminRole sets a user's role to admin
func (r *UserRepository) SetAdminRole(id int64) error {
	query, args, err := r.psql.
		Update("users").
		Set("role", models.RoleAdmin).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// isDuplicateKeyError checks if the error is a duplicate key error
func isDuplicateKeyError(err error) bool {
	// SQLite error for UNIQUE constraint failed
	return err != nil && (err.Error() == "UNIQUE constraint failed: users.email" ||
		err.Error() == "UNIQUE constraint failed: users.google_id")
}
