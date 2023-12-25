package api

import (
	"fmt"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for the banking service.
type Server struct {
    config     util.Config
    store      db.Store
    router     *gin.Engine
    tokenMaker token.Maker
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
    tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
    if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }

    server := &Server{
        config:     config,
        store:      store,
        tokenMaker: tokenMaker,
    }

    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("currency", validCurrency)
    }

    server.setupRouter()
    return server, nil
}

func (server *Server) Start(address string) error {
    return server.router.Run(address)
}

func (server *Server) setupRouter() {
    router := gin.Default()

    router.POST("/users", server.CreateUser)
    router.POST("/users/login", server.LoginUser)
    router.POST("/token/renew_access", server.RenewAccessToken)

    authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

    authRoutes.POST("/accounts", server.CreateAccount)
    authRoutes.GET("/accounts/:id", server.GetAccount)
    authRoutes.GET("/accounts", server.ListAccounts)

    authRoutes.POST("/transfers", server.CreateTransfer)

    server.router = router
}

func errorResponse(err error) gin.H {
    return gin.H{"error": err.Error()}
}
