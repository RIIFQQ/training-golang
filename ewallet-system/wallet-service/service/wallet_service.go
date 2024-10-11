package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"training-golang/ewallet-system/wallet-service/entity"
)

// IWalletService defines the methods for handling wallet operations.
type IWalletService interface {
	CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error)
	Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error)
	GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error)
}

// IWalletRepository defines the methods for data access related to wallets.
type IWalletRepository interface {
	CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error)
	TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error)
	Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error)
	GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error)
}

type WalletServiceImpl struct {
	repo IWalletRepository
}

func NewWalletService(repo IWalletRepository) IWalletService {
	return &WalletServiceImpl{repo: repo}
}

func (s *WalletServiceImpl) CreateWallet(ctx context.Context, userID int32) (*entity.Wallet, error) {
	wallet, err := s.repo.CreateWallet(ctx, userID)
	if err != nil {
		log.Printf("failed to create wallet for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}
	return wallet, nil
}

func (s *WalletServiceImpl) TopUp(ctx context.Context, userID int32, amount float32) (*entity.Wallet, error) {
	wallet, err := s.repo.TopUp(ctx, userID, amount)
	if err != nil {
		log.Printf("failed to top up wallet for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to top up wallet: %w", err)
	}
	return wallet, nil
}

func (s *WalletServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID int32, amount float32) (*entity.Wallet, error) {
	fromWallet, err := s.repo.GetWallet(ctx, fromUserID)
	if err != nil {
		log.Printf("failed to retrieve wallet for user %d: %v", fromUserID, err)
		return nil, fmt.Errorf("failed to retrieve wallet: %w", err)
	}

	if fromWallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	toWallet, err := s.repo.GetWallet(ctx, toUserID)
	if err != nil {
		log.Printf("failed to retrieve wallet for user %d: %v", toUserID, err)
		return nil, fmt.Errorf("failed to retrieve wallet: %w", err)
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	_, err = s.repo.Transfer(ctx, fromUserID, toUserID, amount)
	if err != nil {
		log.Printf("failed to transfer amount from user %d to user %d: %v", fromUserID, toUserID, err)
		return nil, fmt.Errorf("failed to transfer amount: %w", err)
	}

	return fromWallet, nil
}

func (s *WalletServiceImpl) GetWallet(ctx context.Context, userID int32) (*entity.Wallet, error) {
	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve wallet for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve wallet: %w", err)
	}
	return wallet, nil
}

func (s *WalletServiceImpl) GetTransactions(ctx context.Context, userID int32) ([]*entity.Transaction, error) {
	transactions, err := s.repo.GetTransactions(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve transactions for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve transactions: %w", err)
	}
	return transactions, nil
}
