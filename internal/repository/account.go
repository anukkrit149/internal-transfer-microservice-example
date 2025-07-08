package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"internal-transfer-microservice/internal/domain/account"
	"internal-transfer-microservice/internal/infrastructure/db"
)

type AccountRepoImpl struct {
	db db.Database
}

// GetConn Helper to get the DB connection
func (a *AccountRepoImpl) GetConn() *gorm.DB {
	return a.db.GetConnection()
}

func (a *AccountRepoImpl) UpdateAccountsInTx(ctx context.Context, srcAccount *account.Model, destAccount *account.Model) error {
	err := a.GetConn().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(srcAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(destAccount).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountRepoImpl) GetAccount(ctx context.Context, accountId string) (*account.Model, error) {
	var acc account.Model
	err := a.GetConn().First(&acc, "account_id = ?", accountId)
	if err.Error != nil {
		return nil, err.Error
	}
	return &acc, nil
}

func (a *AccountRepoImpl) UpdateAccount(ctx context.Context, account *account.Model) error {
	err := a.GetConn().Save(account)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (a *AccountRepoImpl) CreateAccount(ctx context.Context, accountModel *account.Model) error {
	var tempAccount *account.Model
	err := a.GetConn().First(&tempAccount, "account_id = ?", accountModel.AccountId)
	if err.Error == nil {
		return errors.New("account already exists")
	}
	err = a.GetConn().Create(accountModel)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func NewAccountRepo(db db.Database) *AccountRepoImpl {
	return &AccountRepoImpl{
		db: db,
	}
}
