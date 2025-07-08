package account

import "internal-transfer-microservice/internal/domain"

type Model struct {
	domain.Base
	AccountId string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}
