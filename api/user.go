package api

import (
	"errors"
	"net/http"
	"time"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CreateUserRequest struct {
    Username string `json:"username" binding:"required,alphanum"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"full_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
    Username          string    `json:"username"`
    FullName          string    `json:"full_name"`
    Email             string    `json:"email"`
    PasswordChangedAt time.Time `json:"password_changed_at"`
    CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
    var req CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    hashedPassword, err := util.HashPassword(req.Password)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    arg := db.CreateUserParams{
        Username:       req.Username,
        HashedPassword: hashedPassword,
        FullName:       req.FullName,
        Email:          req.Email,
    }

    user, err := server.store.CreateUser(ctx, arg)
    if err != nil {
        errCode := db.ErrorCode(err)
        if errCode == db.UNIQUE_VIOLATION {
            ctx.JSON(http.StatusForbidden, errorResponse(err))
            return
        }

        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    resp := CreateUserResponse{
        Username:          user.Username,
        FullName:          user.FullName,
        Email:             user.Email,
        PasswordChangedAt: user.PasswordChangedAt,
        CreatedAt:         user.CreatedAt,
    }
    ctx.JSON(http.StatusOK, resp)
}

type LoginUserRequest struct {
    Username string `json:"username" binding:"required,alphanum"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
    AccessToken string             `json:"access_token"`
    User        CreateUserResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
    var req LoginUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    user, err := server.store.GetUser(ctx, req.Username)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    err = util.CheckPassword(req.Password, user.HashedPassword)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    resp := LoginUserResponse{
        AccessToken: accessToken,
        User: CreateUserResponse{
            Username:          user.Username,
            FullName:          user.FullName,
            Email:             user.Email,
            PasswordChangedAt: user.PasswordChangedAt,
            CreatedAt:         user.CreatedAt,
        },
    }
    ctx.JSON(http.StatusOK, resp)
}
