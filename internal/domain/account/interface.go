package account

import "context"

type Repository interface {
	GetAccount(ctx context.Context, accountId string) (*Model, error)
	UpdateAccount(ctx context.Context, account *Model) error
	CreateAccount(ctx context.Context, account *Model) error
	UpdateAccountsInTx(ctx context.Context, srcAccount *Model, destAccount *Model) error
}

type Service interface {
	GetAccount(ctx context.Context, accountId string) (*GetAccountResponse, error)
	CreateAccount(ctx context.Context, accountId string, balance float64) (ApiResponse, error)
	TxnAccount(ctx context.Context, accountId, destinationAccountId string, amount float64) (ApiResponse, error)
}
