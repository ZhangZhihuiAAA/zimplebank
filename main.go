package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ZhangZhihuiAAA/zimplebank/api"
	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	_ "github.com/ZhangZhihuiAAA/zimplebank/doc/statik"
	"github.com/ZhangZhihuiAAA/zimplebank/gapi"
	"github.com/ZhangZhihuiAAA/zimplebank/mail"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/ZhangZhihuiAAA/zimplebank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
    dbMigrationRetryLimit = 6
    dbMigrationRetryInterval = 20 * time.Second
)

func main() {
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal().Err(err).Msg("failed to load config")
    }

    if config.Environment == "DEV" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }

    runDBMigration(config.MigrationURL, config.DBSource, dbMigrationRetryLimit, dbMigrationRetryInterval)

    connPool, err := pgxpool.New(context.Background(), config.DBSource)
    if err != nil {
        log.Fatal().Err(err)
    }

    store := db.NewStore(connPool)

    redisOpt := asynq.RedisClientOpt{
        Addr: config.RedisAddress,
    }

    taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
    go runTaskProcessor(config, redisOpt, store)
    go runGatewayServer(config, store, taskDistributor)
    runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string, retryLimit int, retryInterval time.Duration) {
    retryTimes := 0
RETRY:
    migration, err := migrate.New(migrationURL, dbSource)
    if err != nil {
        if strings.Contains(err.Error(), "connection refused") && retryTimes <= retryLimit {
            log.Info().Msg("retrying creating migrate instance")
            time.Sleep(retryInterval)
            retryTimes++
            goto RETRY
        }
        log.Fatal().Err(err).Msg("failed to create migrate instance")
    }

    if err = migration.Up(); err != nil && err != migrate.ErrNoChange{
        log.Fatal().Err(err).Msg("failed to run migration up")
    }

    log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
    mailer := mail.NewOneTwoSixSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
    taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
    log.Info().Msg("start task processor")
    err := taskProcessor.Start()
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start task processor")
    }
}

func runGinServer(config util.Config, store db.Store) {
    server, err := api.NewServer(config, store)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create api server")
    }

    log.Info().Msg("Start api server")
    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start api server:")
    }
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := gapi.NewServer(config, store, taskDistributor)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create gapi server")
    }

    jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
        MarshalOptions: protojson.MarshalOptions{
            UseProtoNames: true,
        },
        UnmarshalOptions: protojson.UnmarshalOptions{
            DiscardUnknown: true,
        },
    })

    grpcMux := runtime.NewServeMux(jsonOption)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    err = pb.RegisterZimpleBankHandlerServer(ctx, grpcMux, server)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to register handler server")
    }

    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    statikFS, err := fs.New()
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create statik fs:")
    }

    swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
    mux.Handle("/swagger/", swaggerHandler)

    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create listener")
    }

    log.Info().Msgf("Start HTTP gateway server at %s", listener.Addr().String())
    handler := gapi.HttpLogger(mux)
    err = http.Serve(listener, handler)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start HTTP gateway server")
    }
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
    server, err := gapi.NewServer(config, store, taskDistributor)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create gapi server")
    }

    grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
    grpcServer := grpc.NewServer(grpcLogger)
    pb.RegisterZimpleBankServer(grpcServer, server)

    // This step is optional, but is highly recommended.
    // It allows GRPC client easily explore what RPCs are available on the server and how to call them.
    reflection.Register(grpcServer)

    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create listener")
    }

    log.Info().Msgf("Start GRPC server at %s", listener.Addr().String())
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start GRPC server")
    }
}
