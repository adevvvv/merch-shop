package repositories

import (
	"context"
	"database/sql"
	"shop/internal/models"
)

type TransactionRepository interface {
	Create(ctx context.Context, fromUser, toUser, amount int) error
	GetReceived(ctx context.Context, userID int) ([]models.CoinTransactionReceived, error)
	GetSent(ctx context.Context, userID int) ([]models.CoinTransactionSent, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, fromUser, toUser, amount int) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO coin_transactions(from_user, to_user, amount) VALUES($1, $2, $3)",
		fromUser, toUser, amount,
	)
	return err
}

func (r *transactionRepository) GetReceived(ctx context.Context, userID int) ([]models.CoinTransactionReceived, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT u.username, ct.amount, ct.created_at 
		FROM coin_transactions ct
		JOIN users u ON ct.from_user = u.id
		WHERE ct.to_user = $1
		ORDER BY ct.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.CoinTransactionReceived
	for rows.Next() {
		var t models.CoinTransactionReceived
		if err := rows.Scan(&t.FromUser, &t.Amount, &t.Date); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepository) GetSent(ctx context.Context, userID int) ([]models.CoinTransactionSent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT u.username, ct.amount, ct.created_at 
		FROM coin_transactions ct
		JOIN users u ON ct.to_user = u.id
		WHERE ct.from_user = $1
		ORDER BY ct.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.CoinTransactionSent
	for rows.Next() {
		var t models.CoinTransactionSent
		if err := rows.Scan(&t.ToUser, &t.Amount, &t.Date); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}