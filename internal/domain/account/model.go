package account

import "internal-transfer-microservice/internal/domain"

type Model struct {
	domain.Base
	AccountId string  `json:"account_id" gorm:"uniqueIndex;"`
	Balance   float64 `json:"balance"`
}

func (Model) TableName() string {
	return "accounts"
}
