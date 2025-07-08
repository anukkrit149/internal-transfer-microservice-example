package account

type GetAccountResponse struct {
	AccountId string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type ApiResponse struct {
	Message string `json:"message"`
}

type CreateAccountRequest struct {
	AccountId      string  `json:"account_id"`
	InitialBalance float64 `json:"initial_balance"`
}

type TxnAccountRequest struct {
	SourceAccountId      string  `json:"account_id"`
	DestinationAccountId string  `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
}
