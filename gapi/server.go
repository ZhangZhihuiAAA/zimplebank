package gapi

import (
	"fmt"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/ZhangZhihuiAAA/zimplebank/worker"
)

// Server serves GRPC requests for the banking service.
type Server struct {
    pb.UnimplementedZimpleBankServer
    config          util.Config
    store           db.Store
    tokenMaker      token.Maker
    taskDistributor worker.TaskDistributor
}

// NewServer creates a new GRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
    tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
    if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }

    server := &Server{
        config:     config,
        store:      store,
        tokenMaker: tokenMaker,
        taskDistributor: taskDistributor,
    }

    return server, nil
}
