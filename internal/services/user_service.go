// internal/services/user_service.go
package services

import (
	"context"
	"errors"
	"shop/internal/models"
	"shop/internal/repositories"
)

type UserService struct {
	userRepo        repositories.UserRepository
	inventoryRepo   repositories.InventoryRepository
	transactionRepo repositories.TransactionRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	inventoryRepo repositories.InventoryRepository,
	transactionRepo repositories.TransactionRepository,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		inventoryRepo:   inventoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *UserService) GetUserInfo(ctx context.Context, userID int) (models.InfoResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return models.InfoResponse{}, err
	}

	type result struct {
		items interface{}
		err   error
	}

	invChan := make(chan result)
	recvChan := make(chan result)
	sentChan := make(chan result)

	go func() {
		items, err := s.inventoryRepo.GetByUserID(ctx, userID)
		invChan <- result{items: items, err: err}
	}()

	go func() {
		transactions, err := s.transactionRepo.GetReceived(ctx, userID)
		recvChan <- result{items: transactions, err: err}
	}()

	go func() {
		transactions, err := s.transactionRepo.GetSent(ctx, userID)
		sentChan <- result{items: transactions, err: err}
	}()

	var (
		inventory   []models.InventoryItem
		received    []models.CoinTransactionReceived
		sent        []models.CoinTransactionSent
		collectErrs []error
	)

	for i := 0; i < 3; i++ {
		select {
		case res := <-invChan:
			if res.err != nil {
				collectErrs = append(collectErrs, res.err)
			} else {
				inventory = res.items.([]models.InventoryItem)
			}

		case res := <-recvChan:
			if res.err != nil {
				collectErrs = append(collectErrs, res.err)
			} else {
				received = res.items.([]models.CoinTransactionReceived)
			}

		case res := <-sentChan:
			if res.err != nil {
				collectErrs = append(collectErrs, res.err)
			} else {
				sent = res.items.([]models.CoinTransactionSent)
			}
		}
	}

	if len(collectErrs) > 0 {
		return models.InfoResponse{}, errors.New("failed to collect user info")
	}

	return models.InfoResponse{
        Coins:     user.Coins,
        Inventory: inventoryItems,
        CoinHistory: models.CoinHistoryResponse{
            Received: receivedTransactions,
            Sent:     sentTransactions,
        },
    }, nil
}

func (s *UserService) BuyItem(ctx context.Context, userID int, item string) error {
	price, ok := models.MerchItems[item]
	if !ok {
		return errors.New("invalid item")
	}

	coins, err := s.userRepo.GetCoins(ctx, userID)
	if err != nil {
		return err
	}

	if coins < price {
		return errors.New("insufficient funds")
	}

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.userRepo.UpdateCoins(ctx, userID, coins-price); err != nil {
		return err
	}

	if err := s.inventoryRepo.AddItem(ctx, userID, item, 1); err != nil {
		return err
	}

	return tx.Commit()
}