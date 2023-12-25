package api

import (
	"errors"
	"net/http"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
    Owner    string `json:"owner" binding:"required"`
    Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
    var req CreateAccountRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    arg := db.CreateAccountParams{
        Owner:    authPayload.Username,
        Currency: req.Currency,
        Balance:  0.00,
    }

    account, err := server.store.CreateAccount(ctx, arg)
    if err != nil {
        errCode := db.ErrorCode(err)
        if errCode == db.FOREIGN_KEY_VIOLATION || errCode == db.UNIQUE_VIOLATION {
            ctx.JSON(http.StatusForbidden, errorResponse(err))
            return
        }

        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
    var req GetAccountRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    account, err := server.store.GetAccount(ctx, req.ID)
    if err != nil {
        if errors.Is(db.ErrNoRows, err) {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }

        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    if account.Owner != authPayload.Username {
        err := errors.New("account does not belong to the authenticated user")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
    PageID   int32 `form:"page_id" binding:"required,min=1"`
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccounts(ctx *gin.Context) {
    var req ListAccountsRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    arg := db.ListAccountsParams{
        Owner:  authPayload.Username,
        Limit:  req.PageSize,
        Offset: (req.PageID - 1) * req.PageSize,
    }

    accounts, err := server.store.ListAccounts(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, accounts)
}
