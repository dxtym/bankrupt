package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/token"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=8"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	fromAccount, valid := s.validateCurrency(ctx, arg.FromAccountId, req.Currency)
	if !valid {
		return
	} 
	
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("account has limited privileges")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	if _, valid	= s.validateCurrency(ctx, arg.ToAccountId, req.Currency); !valid {
		return
	}

	result, err := s.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) validateCurrency(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] mismatch %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
