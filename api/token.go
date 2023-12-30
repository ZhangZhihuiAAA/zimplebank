package api

import (
	"errors"
	"net/http"
	"time"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type RenewAccessTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
    AccessToken          string    `json:"access_token"`
    AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) RenewAccessToken(ctx *gin.Context) {
    var req RenewAccessTokenRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    session, err := server.store.GetSession(ctx, refreshPayload.ID)
    if err != nil {
        if errors.Is(err, db.ErrNoRows) {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    if session.IsBlocked {
        err := errors.New("session is blocked")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    if session.Username != refreshPayload.Username {
        err := errors.New("incorrect session user")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    if session.RefreshToken != req.RefreshToken {
        err := errors.New("mismatched session token")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    if time.Now().After(session.ExpiresAt) {
        err := errors.New("expired session")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    accessToken, accessPayload, err := server.tokenMaker.CreateToken(
        refreshPayload.Username,
        refreshPayload.Role,
        server.config.AccessTokenDuration,
    )
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    resp := RenewAccessTokenResponse{
        AccessToken:           accessToken,
        AccessTokenExpiresAt:  accessPayload.ExpiresAt,
    }
    ctx.JSON(http.StatusOK, resp)
}
