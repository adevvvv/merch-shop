package repositories

import (
	"context"
	"database/sql"
	"time"
)

type SessionRepository interface {
	Create(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetByUserID(ctx context.Context, userID int) (string, error)
	Delete(ctx context.Context, userID int) error
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO sessions(user_id, token, expires_at) VALUES($1, $2, $3)",
		userID, token, expiresAt,
	)
	return err
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID int) (string, error) {
	var token string
	err := r.db.QueryRowContext(ctx,
		"SELECT token FROM sessions WHERE user_id = $1",
		userID,
	).Scan(&token)
	return token, err
}

func (r *sessionRepository) Delete(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM sessions WHERE user_id = $1",
		userID,
	)
	return err
}