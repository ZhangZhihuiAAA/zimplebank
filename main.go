package main

import (
	"context"
	"log"

	"github.com/ZhangZhihuiAAA/zimplebank/api"
	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }

    connPool, err := pgxpool.New(context.Background(), config.DBSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    store := db.NewStore(connPool)
    server, err := api.NewServer(config, store)
    if err != nil {
        log.Fatal("cannot create server:", err)
    }

    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal("cannot start server:", err)
    }
}