package service

import (
	"context"
	"errors"
	"fmt"
	"internal-transfer-microservice/internal/domain/account"
	"internal-transfer-microservice/internal/infrastructure/cache"
	"time"
)

const UpdateAccountResourceLockKey = "update_account:%s"

var ErrAccountNotFound = errors.New("account not found")

type AccountServiceImpl struct {
	cache            cache.Cache
	repo             account.Repository
	lockPollInterval time.Duration
}

func (a *AccountServiceImpl) acquireLockWithPolling(ctx context.Context, key string, lockTTL, waitTimeout time.Duration) (token string, err error) {
	deadline := time.Now().Add(waitTimeout)
	for {
		_, err = a.cache.Lock(ctx, key, 0)
		if err == nil {
			return key, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("timed out waiting for %s", key)
		}
		time.Sleep(a.lockPollInterval)
	}
}

func (a *AccountServiceImpl) GetAccount(ctx context.Context, accountId string) (*account.GetAccountResponse, error) {
	acc, err := a.repo.GetAccount(ctx, accountId)
	if err != nil {
		return nil, err
	}
	response := &account.GetAccountResponse{
		AccountId: acc.AccountId,
		Balance:   acc.Balance,
	}

	return response, nil
}

func (a *AccountServiceImpl) CreateAccount(ctx context.Context, accountId string, balance float64) (account.ApiResponse, error) {
	newAccount := &account.Model{
		AccountId: accountId,
		Balance:   balance,
	}

	err := a.repo.CreateAccount(ctx, newAccount)
	if err != nil {
		return account.ApiResponse{Message: "Failed to create account"}, err
	}

	return account.ApiResponse{Message: "Account created successfully"}, nil
}

func (a *AccountServiceImpl) TxnAccount(ctx context.Context, sourceAccountId, destAccountId string, amount float64) (account.ApiResponse, error) {
	lock1Key := fmt.Sprintf(UpdateAccountResourceLockKey, sourceAccountId)
	lock2Key := fmt.Sprintf(UpdateAccountResourceLockKey, destAccountId)
	if sourceAccountId > destAccountId {
		lock1Key, lock2Key = lock2Key, lock1Key // to avoid deadlock, in case of concurrent txn (A->B) & (B->A)
	}

	// resource locking
	resource1, err := a.acquireLockWithPolling(ctx, lock1Key, 0, 100*time.Millisecond)
	if err != nil {
		return account.ApiResponse{Message: "Failed to acquire lock for transaction"}, err
	}
	defer a.cache.Release(ctx, resource1)

	resource2, err := a.acquireLockWithPolling(ctx, lock2Key, 0, 100*time.Millisecond)
	if err != nil {
		return account.ApiResponse{Message: "Failed to acquire lock for transaction"}, err
	}
	defer a.cache.Release(ctx, resource2)

	sourceAccount, err := a.repo.GetAccount(ctx, sourceAccountId)
	if err != nil {
		return account.ApiResponse{Message: "Source account not found"}, ErrAccountNotFound
	}

	if sourceAccount.Balance < amount {
		return account.ApiResponse{Message: "Insufficient balance"}, nil
	}

	destAccount, err := a.repo.GetAccount(ctx, destAccountId)
	if err != nil {
		return account.ApiResponse{Message: "Destination account not found"}, ErrAccountNotFound
	}

	sourceAccount.Balance -= amount
	destAccount.Balance += amount
	err = a.repo.UpdateAccountsInTx(ctx, sourceAccount, destAccount)
	if err != nil {
		return account.ApiResponse{Message: "Transaction failed during database update"}, err
	}

	return account.ApiResponse{Message: "Transaction completed successfully"}, nil
}

func NewAccountService(repo account.Repository, cache cache.Cache) account.Service {
	return &AccountServiceImpl{
		repo:             repo,
		cache:            cache,
		lockPollInterval: 10 * time.Millisecond,
	}
}
