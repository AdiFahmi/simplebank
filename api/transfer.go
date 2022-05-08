package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/adifahmi/simplebank/db/sqlc"
	"github.com/adifahmi/simplebank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	ToAccountID int64  `json:"to_account_id" binding:"required,min=1"`
	Ammount     int64  `json:"ammount" binding:"required,gt=0"`
	Currency    string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// check if from account exists
	_, valid := server.validAccount(ctx, authPayload.UserID, req.Currency)
	if !valid {
		return
	}

	// check if to account exists
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: authPayload.UserID,
		ToAccountID:   req.ToAccountID,
		Ammount:       req.Ammount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			err := fmt.Errorf("account %d not found", accountID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
