package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ZhangZhihuiAAA/zimplebank/api"
	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
    RETRY_INTERVAL = 5 * time.Second
    RETRY_LIMIT = 10
)

var retryTimes = 0
var config util.Config
var connPool *pgxpool.Pool
var ctx = context.Background()
var err error

func init() {
    config, err = util.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }

    connPool, err = pgxpool.New(ctx, config.DBSource)
    if err != nil {
        log.Fatal(err)
    }

CONNECT_DB:
    _, err = connPool.Query(ctx, "SELECT * FROM users LIMIT 1;")
    if err != nil {
        if strings.Contains(err.Error(), "failed to connect to") {
            if retryTimes < RETRY_LIMIT {
                retryTimes++
                time.Sleep(RETRY_INTERVAL)
                goto CONNECT_DB
            } else {
                log.Fatal("failed to connect to db:", err)
            }
        }
    
        if strings.Contains(err.Error(), "does not exist") {
            data, err := os.ReadFile(config.DBInitSchemaFile)
            if err != nil {
                log.Fatalf("failed to read file: %s", config.DBInitSchemaFile)
            }
            connPool.Exec(ctx, string(data))
            if err != nil {
                log.Fatal("error occurred when init schema:", err)
            }
        } else {
            log.Fatal("error occurred:", err)
        }
    }
}

func main() {
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