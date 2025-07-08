package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"internal-transfer-microservice/internal/domain/account"
)

type AccountController struct {
	accountService account.Service
}

// NewAccountController creates a new AccountController
func NewAccountController(accountService account.Service) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// GetAccount handles GET /accounts/:id
func (c *AccountController) GetAccount(ctx *gin.Context) {
	accountId := ctx.Param("id")

	response, err := c.accountService.GetAccount(ctx, accountId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateAccount handles POST /accounts
func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req account.CreateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.accountService.CreateAccount(ctx, req.AccountId, req.InitialBalance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// TransferMoney handles POST /accounts/transfer
func (c *AccountController) TransferMoney(ctx *gin.Context) {
	var req account.TxnAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.accountService.TxnAccount(ctx, req.SourceAccountId, req.DestinationAccountId, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
