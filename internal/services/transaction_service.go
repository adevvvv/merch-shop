// internal/services/transaction_service.go
package services

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/models"
	"shop/internal/repositories"
)

type TransactionService struct {
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(
	userRepo repositories.UserRepository,
	transactionRepo repositories.TransactionRepository,
) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) SendCoin(ctx context.Context, senderID int, req models.SendCoinRequest) error {
	// Находим получателя
	receiver, err := s.userRepo.GetByUsername(ctx, req.ToUser)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("receiver not found")
	}
	if err != nil {
		return err
	}

	// Проверяем баланс отправителя
	senderCoins, err := s.userRepo.GetCoins(ctx, senderID)
	if err != nil {
		return err
	}
	if senderCoins < req.Amount {
		return errors.New("insufficient funds")
	}

	// Выполняем транзакцию
	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем балансы
	if err := s.userRepo.UpdateCoins(ctx, senderID, senderCoins-req.Amount); err != nil {
		return err
	}
	if err := s.userRepo.UpdateCoins(ctx, receiver.ID, receiver.Coins+req.Amount); err != nil {
		return err
	}

	// Записываем транзакцию
	if err := s.transactionRepo.Create(ctx, senderID, receiver.ID, req.Amount); err != nil {
		return err
	}

	return tx.Commit()
}