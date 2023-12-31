package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
    FromAccountID int64   `json:"from_account_id" binding:"required,min=1"`
    ToAccountID   int64   `json:"to_account_id" binding:"required,min=1"`
    Amount        float64 `json:"amount" binding:"required,gt=0.00"`
    Currency      string  `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
    var req TransferRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
    if !valid {
        return
    }

    authPayload := ctx.MustGet(AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)
    if fromAccount.Owner != authPayload.Username {
        err := errors.New("from account does not belong to the authenticated user")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    if _, valid := server.validAccount(ctx, req.ToAccountID, req.Currency); !valid {
        return
    }

    arg := db.TransferTxParams{
        FromAccountID: req.FromAccountID,
        ToAccountID:   req.ToAccountID,
        Amount:        req.Amount,
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
        if errors.Is(db.ErrNoRows, err) {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return account, false
        }

        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return account, false
    }

    if account.Currency != currency {
        err := fmt.Errorf("account [%d] currency mismatch: account currency %s vs transfer currency %s", account.ID, account.Currency, currency)
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return account, false
    }

    return account, true
}