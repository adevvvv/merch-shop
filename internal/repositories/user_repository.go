package repositories

import (
	"context"
	"database/sql"
	"shop/internal/models"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	UpdateCoins(ctx context.Context, userID int, coins int) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetCoins(ctx context.Context, userID int) (int, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, password, coins, is_admin FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Coins, &user.IsAdmin)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users(username, password, coins, is_admin) VALUES($1, $2, $3, $4) RETURNING id",
		user.Username, user.Password, user.Coins, user.IsAdmin,
	).Scan(&user.ID)
	return err
}

func (r *userRepository) UpdateCoins(ctx context.Context, userID int, coins int) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET coins = $1 WHERE id = $2",
		coins, userID,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, username, password, coins, is_admin FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Coins, &user.IsAdmin)
	return user, err
}

func (r *userRepository) GetCoins(ctx context.Context, userID int) (int, error) {
	var coins int
	err := r.db.QueryRowContext(ctx,
		"SELECT coins FROM users WHERE id = $1",
		userID,
	).Scan(&coins)
	return coins, err
}