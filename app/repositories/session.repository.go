package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/maulanashalihin/laju-go/app/models"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type SessionRepository struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

// Create creates a new session
func (r *SessionRepository) Create(session *models.Session) error {
	query, args, err := r.psql.
		Insert("sessions").
		Columns("id", "user_id", "data", "expires_at", "created_at", "updated_at").
		Values(session.ID, session.UserID, session.Data, session.ExpiresAt, time.Now(), time.Now()).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetByID finds a session by ID
func (r *SessionRepository) GetByID(id string) (*models.Session, error) {
	query, args, err := r.psql.
		Select("id", "user_id", "data", "expires_at", "created_at", "updated_at").
		From("sessions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	session := &models.Session{}
	err = r.db.QueryRow(query, args...).Scan(
		&session.ID, &session.UserID, &session.Data,
		&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		// Delete expired session
		r.Delete(id)
		return nil, ErrSessionExpired
	}

	return session, nil
}

// GetByUserID finds all sessions for a user
func (r *SessionRepository) GetByUserID(userID int64) ([]*models.Session, error) {
	query, args, err := r.psql.
		Select("id", "user_id", "data", "expires_at", "created_at", "updated_at").
		From("sessions").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.UserID, &session.Data,
			&session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Skip expired sessions
		if session.ExpiresAt.After(time.Now()) {
			sessions = append(sessions, session)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Update updates an existing session
func (r *SessionRepository) Update(session *models.Session) error {
	query, args, err := r.psql.
		Update("sessions").
		Set("data", session.Data).
		Set("expires_at", session.ExpiresAt).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": session.ID}).
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
		return ErrSessionNotFound
	}

	return nil
}

// Delete deletes a session by ID
func (r *SessionRepository) Delete(id string) error {
	query, args, err := r.psql.
		Delete("sessions").
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
		return ErrSessionNotFound
	}

	return nil
}

// DeleteByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(userID int64) error {
	query, args, err := r.psql.
		Delete("sessions").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteExpired deletes all expired sessions
func (r *SessionRepository) DeleteExpired() error {
	query, args, err := r.psql.
		Delete("sessions").
		Where(squirrel.Lt{"expires_at": time.Now()}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// DecodeData decodes session data into SessionData struct
func (r *SessionRepository) DecodeData(data string) (*models.SessionData, error) {
	var sessionData models.SessionData
	err := json.Unmarshal([]byte(data), &sessionData)
	if err != nil {
		return nil, err
	}
	return &sessionData, nil
}

// EncodeData encodes SessionData into JSON string
func (r *SessionRepository) EncodeData(data *models.SessionData) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
