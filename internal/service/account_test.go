package service

import (
	"context"
	"errors"
	"internal-transfer-microservice/internal/domain/account"
	"sync"
	"testing"
	"time"
)

// MockRepository is a mock implementation of account.Repository
type MockRepository struct {
	accounts map[string]*account.Model
	mu       sync.Mutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		accounts: make(map[string]*account.Model),
	}
}

func (m *MockRepository) GetAccount(ctx context.Context, accountId string) (*account.Model, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	acc, exists := m.accounts[accountId]
	if !exists {
		return nil, errors.New("account not found")
	}
	return acc, nil
}

func (m *MockRepository) UpdateAccount(ctx context.Context, account *account.Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.accounts[account.AccountId]
	if !exists {
		return errors.New("account not found")
	}
	m.accounts[account.AccountId] = account
	return nil
}

func (m *MockRepository) CreateAccount(ctx context.Context, account *account.Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.accounts[account.AccountId]
	if exists {
		return errors.New("account already exists")
	}
	m.accounts[account.AccountId] = account
	return nil
}

func (m *MockRepository) UpdateAccountsInTx(ctx context.Context, srcAccount *account.Model, destAccount *account.Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, srcExists := m.accounts[srcAccount.AccountId]
	if !srcExists {
		return errors.New("source account not found")
	}

	_, destExists := m.accounts[destAccount.AccountId]
	if !destExists {
		return errors.New("destination account not found")
	}

	m.accounts[srcAccount.AccountId] = srcAccount
	m.accounts[destAccount.AccountId] = destAccount
	return nil
}

// MockCache is a mock implementation of cache.Cache
type MockCache struct {
	locks map[string]bool
	data  map[string]string
	mu    sync.Mutex
}

func NewMockCache() *MockCache {
	return &MockCache{
		locks: make(map[string]bool),
		data:  make(map[string]string),
	}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, exists := m.data[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (m *MockCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}

func (m *MockCache) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.locks[key] {
		return false, errors.New("lock already acquired")
	}
	m.locks[key] = true
	return true, nil
}

func (m *MockCache) Release(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.locks, key)
	return nil
}

func (m *MockCache) Close() error {
	return nil
}

// Test cases
func TestCreateAccount(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	cache := NewMockCache()
	service := NewAccountService(repo, cache)
	ctx := context.Background()

	// Test case: Create a new account
	accountId := "acc123"
	balance := 1000.0

	response, err := service.CreateAccount(ctx, accountId, balance)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Message != "Account created successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}

	// Verify account was created
	account, err := repo.GetAccount(ctx, accountId)
	if err != nil {
		t.Errorf("Expected account to exist, got error: %v", err)
	}
	if account.AccountId != accountId {
		t.Errorf("Expected account ID %s, got %s", accountId, account.AccountId)
	}
	if account.Balance != balance {
		t.Errorf("Expected balance %.2f, got %.2f", balance, account.Balance)
	}
}

func TestGetExistingAccount(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	cache := NewMockCache()
	service := NewAccountService(repo, cache)
	ctx := context.Background()

	// Create an account first
	accountId := "acc123"
	balance := 1000.0
	repo.CreateAccount(ctx, &account.Model{
		AccountId: accountId,
		Balance:   balance,
	})

	// Test case: Get an existing account
	response, err := service.GetAccount(ctx, accountId)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.AccountId != accountId {
		t.Errorf("Expected account ID %s, got %s", accountId, response.AccountId)
	}
	if response.Balance != balance {
		t.Errorf("Expected balance %.2f, got %.2f", balance, response.Balance)
	}
}

func TestGetNonExistentAccount(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	cache := NewMockCache()
	service := NewAccountService(repo, cache)
	ctx := context.Background()

	// Test case: Get a non-existent account
	accountId := "non_existent_account"
	_, err := service.GetAccount(ctx, accountId)
	if err == nil {
		t.Error("Expected error for non-existent account, got nil")
	}
}

func TestSimpleTransfer(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	cache := NewMockCache()
	service := NewAccountService(repo, cache)
	ctx := context.Background()

	// Create source and destination accounts
	sourceId := "source123"
	destId := "dest456"
	sourceBalance := 1000.0
	destBalance := 500.0
	transferAmount := 200.0

	repo.CreateAccount(ctx, &account.Model{
		AccountId: sourceId,
		Balance:   sourceBalance,
	})
	repo.CreateAccount(ctx, &account.Model{
		AccountId: destId,
		Balance:   destBalance,
	})

	// Test case: Simple transfer
	response, err := service.TxnAccount(ctx, sourceId, destId, transferAmount)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Message != "Transaction completed successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}

	// Verify balances
	sourceAccount, _ := repo.GetAccount(ctx, sourceId)
	destAccount, _ := repo.GetAccount(ctx, destId)

	expectedSourceBalance := sourceBalance - transferAmount
	expectedDestBalance := destBalance + transferAmount

	if sourceAccount.Balance != expectedSourceBalance {
		t.Errorf("Expected source balance %.2f, got %.2f", expectedSourceBalance, sourceAccount.Balance)
	}
	if destAccount.Balance != expectedDestBalance {
		t.Errorf("Expected destination balance %.2f, got %.2f", expectedDestBalance, destAccount.Balance)
	}
}

func TestConcurrentTransfersOnDifferentAccounts(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	cache := NewMockCache()
	service := NewAccountService(repo, cache)
	ctx := context.Background()

	// Create accounts
	acc1 := "acc1"
	acc2 := "acc2"
	acc3 := "acc3"
	acc4 := "acc4"

	repo.CreateAccount(ctx, &account.Model{AccountId: acc1, Balance: 1000.0})
	repo.CreateAccount(ctx, &account.Model{AccountId: acc2, Balance: 1000.0})
	repo.CreateAccount(ctx, &account.Model{AccountId: acc3, Balance: 1000.0})
	repo.CreateAccount(ctx, &account.Model{AccountId: acc4, Balance: 1000.0})

	// Test case: Concurrent transfers on different accounts
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := service.TxnAccount(ctx, acc1, acc2, 200.0)
		if err != nil {
			t.Errorf("Transfer 1 failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		_, err := service.TxnAccount(ctx, acc3, acc4, 300.0)
		if err != nil {
			t.Errorf("Transfer 2 failed: %v", err)
		}
	}()

	wg.Wait()

	// Verify balances
	acc1Final, _ := repo.GetAccount(ctx, acc1)
	acc2Final, _ := repo.GetAccount(ctx, acc2)
	acc3Final, _ := repo.GetAccount(ctx, acc3)
	acc4Final, _ := repo.GetAccount(ctx, acc4)

	if acc1Final.Balance != 800.0 {
		t.Errorf("Expected acc1 balance 800.0, got %.2f", acc1Final.Balance)
	}
	if acc2Final.Balance != 1200.0 {
		t.Errorf("Expected acc2 balance 1200.0, got %.2f", acc2Final.Balance)
	}
	if acc3Final.Balance != 700.0 {
		t.Errorf("Expected acc3 balance 700.0, got %.2f", acc3Final.Balance)
	}
	if acc4Final.Balance != 1300.0 {
		t.Errorf("Expected acc4 balance 1300.0, got %.2f", acc4Final.Balance)
	}
}

func TestDeadlockPrevention(t *testing.T) {
	// Setup
	repo := NewMockRepository()
	cache := NewMockCache()
	service := NewAccountService(repo, cache)
	ctx := context.Background()

	// Create accounts
	accA := "accA"
	accB := "accB"

	repo.CreateAccount(ctx, &account.Model{AccountId: accA, Balance: 1000.0})
	repo.CreateAccount(ctx, &account.Model{AccountId: accB, Balance: 1000.0})

	// Test case: Concurrent transfers A->B and B->A
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := service.TxnAccount(ctx, accA, accB, 200.0)
		if err != nil {
			t.Errorf("Transfer A->B failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		_, err := service.TxnAccount(ctx, accB, accA, 300.0)
		if err != nil {
			t.Errorf("Transfer B->A failed: %v", err)
		}
	}()

	wg.Wait()

	// Verify final balances
	accAFinal, _ := repo.GetAccount(ctx, accA)
	accBFinal, _ := repo.GetAccount(ctx, accB)

	// The final balances should be consistent regardless of execution order
	// A starts with 1000, loses 200, gains 300 = 1100
	// B starts with 1000, loses 300, gains 200 = 900
	if accAFinal.Balance != 1100.0 {
		t.Errorf("Expected accA balance 1100.0, got %.2f", accAFinal.Balance)
	}
	if accBFinal.Balance != 900.0 {
		t.Errorf("Expected accB balance 900.0, got %.2f", accBFinal.Balance)
	}
}
