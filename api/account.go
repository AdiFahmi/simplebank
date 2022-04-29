package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/VividCortex/mysqlerr"
	db "github.com/adifahmi/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	res, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if sqlErr, ok := err.(*mysql.MySQLError); ok {
			log.Printf("Err %s with number %d", err.Error(), sqlErr.Number)
			if sqlErr.Number == mysqlerr.ER_NO_REFERENCED_ROW_2 {
				ctx.JSON(http.StatusBadRequest, errorStringResponse("Owner doesn't exist"))
				return
			} else if sqlErr.Number == mysqlerr.ER_DUP_ENTRY {
				ctx.JSON(http.StatusBadRequest, errorStringResponse("Account with that currency already exists"))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	createdAccID, err := res.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, createdAccID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccuntRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccuntRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID int32 `form:"page"`
}

var limitPage int32 = 10

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  limitPage,
		Offset: (req.PageID - 1) * limitPage,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
